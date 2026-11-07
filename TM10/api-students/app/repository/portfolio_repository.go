package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"api-students/app/model"
)

var (
	ErrInvalidProjectID = errors.New("format id proyek tidak valid")
	ErrProjectLimit     = errors.New("jumlah proyek sudah mencapai batas")
)

// MaxProjectsPerPortfolio membatasi pertumbuhan array projects, sama seperti
// batas replies di Modul 10 Bagian B.
const MaxProjectsPerPortfolio = 20

type PortfolioRepository interface {
	EnsureIndexes(ctx context.Context) error
	GetByStudentID(ctx context.Context, studentID int) (model.Portfolio, error)
	Upsert(ctx context.Context, studentID int, req model.PutPortfolioRequest) (model.Portfolio, error)
	AddProject(ctx context.Context, studentID int, project model.Project) (model.Portfolio, error)
	RemoveProject(ctx context.Context, studentID int, projectID string) (model.Portfolio, error)
	TechnologyStats(ctx context.Context) ([]model.TechCount, error)
}

// portfolioDocument adalah bentuk portofolio SEBAGAIMANA TERSIMPAN di
// MongoDB. Tag bson terpisah dari tag json milik model, sama seperti pola
// reviewDocument di Bagian B.
type portfolioDocument struct {
	ID           bson.ObjectID         `bson:"_id,omitempty"`
	StudentID    int                   `bson:"student_id"`
	Headline     string                `bson:"headline"`
	Bio          string                `bson:"bio"`
	Projects     []projectDocument     `bson:"projects"`
	Certificates []certificateDocument `bson:"certificates"`
	CreatedAt    time.Time             `bson:"created_at"`
	UpdatedAt    time.Time             `bson:"updated_at"`
}

// portfolioWrite dipakai saat menulis (insert/replace) — TANPA field _id,
// supaya MongoDB mempertahankan _id dokumen yang sudah ada saat upsert.
type portfolioWrite struct {
	StudentID    int                   `bson:"student_id"`
	Headline     string                `bson:"headline"`
	Bio          string                `bson:"bio"`
	Projects     []projectDocument     `bson:"projects"`
	Certificates []certificateDocument `bson:"certificates"`
	CreatedAt    time.Time             `bson:"created_at"`
	UpdatedAt    time.Time             `bson:"updated_at"`
}

type projectDocument struct {
	ID           bson.ObjectID `bson:"_id"`
	Title        string        `bson:"title"`
	Technologies []string      `bson:"technologies"`
	URL          string        `bson:"url"`
	Year         int           `bson:"year"`
}

type certificateDocument struct {
	Name   string `bson:"name"`
	Issuer string `bson:"issuer"`
	Date   string `bson:"date"`
}

func (d portfolioDocument) toModel() model.Portfolio {
	projects := make([]model.Project, 0, len(d.Projects))
	for _, p := range d.Projects {
		techs := p.Technologies
		if techs == nil {
			techs = []string{}
		}
		projects = append(projects, model.Project{
			ID: p.ID.Hex(), Title: p.Title, Technologies: techs,
			URL: p.URL, Year: p.Year,
		})
	}
	certificates := make([]model.Certificate, 0, len(d.Certificates))
	for _, cert := range d.Certificates {
		certificates = append(certificates, model.Certificate{
			Name: cert.Name, Issuer: cert.Issuer, Date: cert.Date,
		})
	}
	return model.Portfolio{
		ID: d.ID.Hex(), StudentID: d.StudentID, Headline: d.Headline, Bio: d.Bio,
		Projects: projects, Certificates: certificates,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type portfolioMongoRepository struct{ coll *mongo.Collection }

func NewPortfolioRepository(db *mongo.Database) PortfolioRepository {
	return &portfolioMongoRepository{coll: db.Collection("student_portfolios")}
}

func portfolioNow() time.Time { return time.Now().UTC().Truncate(time.Millisecond) }

func parseProjectID(id string) (bson.ObjectID, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, ErrInvalidProjectID
	}
	return oid, nil
}

// EnsureIndexes dijalankan sekali saat aplikasi menyala.
func (r *portfolioMongoRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			// Satu mahasiswa hanya boleh memiliki SATU dokumen portofolio.
			// Aturan ini dijaga database, bukan hanya diperiksa di kode.
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("student_unique"),
		},
	})
	if err != nil {
		return fmt.Errorf("membuat index student_portfolios: %w", err)
	}
	return nil
}

func (r *portfolioMongoRepository) GetByStudentID(ctx context.Context, studentID int) (model.Portfolio, error) {
	var doc portfolioDocument
	err := r.coll.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Portfolio{}, ErrNotFound
		}
		return model.Portfolio{}, fmt.Errorf("mengambil portofolio: %w", err)
	}
	return doc.toModel(), nil
}

// Upsert membuat portofolio bila belum ada, atau mengganti seluruhnya bila
// sudah ada — sesuai semantik PUT. created_at dipertahankan dari dokumen
// lama bila ada; dibuat baru hanya saat insert pertama.
func (r *portfolioMongoRepository) Upsert(ctx context.Context, studentID int, req model.PutPortfolioRequest) (model.Portfolio, error) {
	createdAt := portfolioNow()
	var existing portfolioDocument
	if err := r.coll.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&existing); err == nil {
		createdAt = existing.CreatedAt
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}

	projects := make([]projectDocument, 0, len(req.Projects))
	for _, p := range req.Projects {
		techs := p.Technologies
		if techs == nil {
			techs = []string{}
		}
		projects = append(projects, projectDocument{
			ID: bson.NewObjectID(), Title: p.Title, Technologies: techs,
			URL: p.URL, Year: p.Year,
		})
	}
	certificates := make([]certificateDocument, 0, len(req.Certificates))
	for _, c := range req.Certificates {
		certificates = append(certificates, certificateDocument{
			Name: c.Name, Issuer: c.Issuer, Date: c.Date,
		})
	}

	write := portfolioWrite{
		StudentID: studentID, Headline: req.Headline, Bio: req.Bio,
		Projects: projects, Certificates: certificates,
		CreatedAt: createdAt, UpdatedAt: portfolioNow(),
	}

	var doc portfolioDocument
	opts := options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After)
	err := r.coll.FindOneAndReplace(ctx, bson.M{"student_id": studentID}, write, opts).Decode(&doc)
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("menyimpan portofolio: %w", err)
	}
	return doc.toModel(), nil
}

// AddProject menambahkan satu proyek ke array projects.
//
// Syarat "projects.19 tidak ada" berarti array belum berisi 20 elemen
// (indeks 0-19). Pemeriksaan batas dan penambahan terjadi dalam satu
// operasi, sama seperti pola replies.49 di Bagian B.
func (r *portfolioMongoRepository) AddProject(ctx context.Context, studentID int, project model.Project) (model.Portfolio, error) {
	techs := project.Technologies
	if techs == nil {
		techs = []string{}
	}
	doc := projectDocument{
		ID: bson.NewObjectID(), Title: project.Title, Technologies: techs,
		URL: project.URL, Year: project.Year,
	}

	limitField := fmt.Sprintf("projects.%d", MaxProjectsPerPortfolio-1)

	var updated portfolioDocument
	err := r.coll.FindOneAndUpdate(ctx,
		bson.M{"student_id": studentID, limitField: bson.M{"$exists": false}},
		bson.M{"$push": bson.M{"projects": doc}, "$set": bson.M{"updated_at": portfolioNow()}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)

	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menambah proyek: %w", err)
	}

	count, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}
	if count == 0 {
		return model.Portfolio{}, ErrNotFound
	}
	return model.Portfolio{}, ErrProjectLimit
}

func (r *portfolioMongoRepository) RemoveProject(ctx context.Context, studentID int, projectID string) (model.Portfolio, error) {
	oid, err := parseProjectID(projectID)
	if err != nil {
		return model.Portfolio{}, err
	}

	var updated portfolioDocument
	err = r.coll.FindOneAndUpdate(ctx,
		bson.M{"student_id": studentID},
		bson.M{
			"$pull": bson.M{"projects": bson.M{"_id": oid}},
			"$set":  bson.M{"updated_at": portfolioNow()},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Portfolio{}, ErrNotFound
		}
		return model.Portfolio{}, fmt.Errorf("menghapus proyek: %w", err)
	}
	return updated.toModel(), nil
}

// TechnologyStats meringkas 10 teknologi paling banyak dipakai di SELURUH
// portofolio. Dua kali $unwind dibutuhkan: sekali untuk membongkar array
// projects, sekali lagi untuk membongkar array technologies di dalam
// masing-masing proyek.
func (r *portfolioMongoRepository) TechnologyStats(ctx context.Context) ([]model.TechCount, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$projects"}},
		{{Key: "$unwind", Value: "$projects.technologies"}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$projects.technologies",
			"total": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}, {Key: "_id", Value: 1}}}},
		{{Key: "$limit", Value: 10}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("meringkas teknologi: %w", err)
	}
	defer cursor.Close(ctx)

	var out []struct {
		Technology string `bson:"_id"`
		Total      int    `bson:"total"`
	}
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("membaca ringkasan teknologi: %w", err)
	}

	result := make([]model.TechCount, 0, len(out))
	for _, o := range out {
		result = append(result, model.TechCount{Technology: o.Technology, Total: o.Total})
	}
	return result, nil
}
