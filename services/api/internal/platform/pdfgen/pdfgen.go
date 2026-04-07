package pdfgen

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
)

// DietPDFData contains data needed to generate a diet PDF.
type DietPDFData struct {
	PatientName      string
	ProfessionalName string
	DietTitle        string
	Objective        string
	ValidFrom        *time.Time
	ValidUntil       *time.Time
	Meals            []MealData
}

// MealData represents a meal within the diet PDF.
type MealData struct {
	Name  string
	Order int
	Notes string
	Items []MealItemData
}

// MealItemData represents a food item within a meal.
type MealItemData struct {
	FoodName      string
	Quantity      string
	Preparation   string
	Substitutions []SubstitutionData
}

// SubstitutionData represents a food substitution.
type SubstitutionData struct {
	FoodName string
	Quantity string
	Notes    string
}

// DocumentPDFData contains data needed to generate a clinical document PDF.
type DocumentPDFData struct {
	PatientName      string
	ProfessionalName string
	DocumentType     string
	Title            string
	Content          map[string]any
	CreatedAt        time.Time
}

// Generator defines the PDF generation interface.
type Generator interface {
	GenerateDietPDF(ctx context.Context, data DietPDFData) ([]byte, error)
	GenerateDocumentPDF(ctx context.Context, data DocumentPDFData) ([]byte, error)
}

// FPDFGenerator implements Generator using go-pdf/fpdf.
type FPDFGenerator struct{}

// New creates a new FPDFGenerator.
func New() *FPDFGenerator {
	return &FPDFGenerator{}
}

// GenerateDietPDF generates a PDF for a diet plan.
func (g *FPDFGenerator) GenerateDietPDF(ctx context.Context, data DietPDFData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	// Header
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, "Plano Alimentar")
	pdf.Ln(12)

	// Patient/Professional info
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Paciente: %s", data.PatientName))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Profissional: %s", data.ProfessionalName))
	pdf.Ln(6)

	if data.DietTitle != "" {
		pdf.Cell(0, 6, fmt.Sprintf("Dieta: %s", data.DietTitle))
		pdf.Ln(6)
	}
	if data.Objective != "" {
		pdf.Cell(0, 6, fmt.Sprintf("Objetivo: %s", data.Objective))
		pdf.Ln(6)
	}

	pdf.Ln(6)

	// Meals
	for _, meal := range data.Meals {
		pdf.SetFont("Helvetica", "B", 12)
		pdf.Cell(0, 8, meal.Name)
		pdf.Ln(8)

		if meal.Notes != "" {
			pdf.SetFont("Helvetica", "I", 9)
			pdf.MultiCell(0, 5, meal.Notes, "", "", false)
			pdf.Ln(2)
		}

		for _, item := range meal.Items {
			pdf.SetFont("Helvetica", "", 10)
			line := fmt.Sprintf("  - %s (%s)", item.FoodName, item.Quantity)
			if item.Preparation != "" {
				line += " — " + item.Preparation
			}
			pdf.MultiCell(0, 5, line, "", "", false)

			for _, sub := range item.Substitutions {
				pdf.SetFont("Helvetica", "", 9)
				subLine := fmt.Sprintf("      Substituir por: %s (%s)", sub.FoodName, sub.Quantity)
				if sub.Notes != "" {
					subLine += " — " + sub.Notes
				}
				pdf.MultiCell(0, 5, subLine, "", "", false)
			}
		}
		pdf.Ln(4)
	}

	// Footer with date
	pdf.SetFont("Helvetica", "I", 8)
	pdf.Cell(0, 6, fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04")))

	if err := pdf.Error(); err != nil {
		return nil, fmt.Errorf("pdfgen: diet pdf error: %w", err)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdfgen: diet pdf output: %w", err)
	}
	return buf.Bytes(), nil
}

// GenerateDocumentPDF generates a PDF for a clinical document.
func (g *FPDFGenerator) GenerateDocumentPDF(ctx context.Context, data DocumentPDFData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	// Header
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, data.Title)
	pdf.Ln(12)

	// Meta
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Tipo: %s", data.DocumentType))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Paciente: %s", data.PatientName))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Profissional: %s", data.ProfessionalName))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Data: %s", data.CreatedAt.Format("02/01/2006")))
	pdf.Ln(10)

	// Content — render key-value pairs from content map
	pdf.SetFont("Helvetica", "", 10)
	for key, val := range data.Content {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.Cell(0, 6, key+":")
		pdf.Ln(6)
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(0, 5, fmt.Sprintf("%v", val), "", "", false)
		pdf.Ln(3)
	}

	// Footer
	pdf.SetFont("Helvetica", "I", 8)
	pdf.Cell(0, 6, fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04")))

	if err := pdf.Error(); err != nil {
		return nil, fmt.Errorf("pdfgen: document pdf error: %w", err)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdfgen: document pdf output: %w", err)
	}
	return buf.Bytes(), nil
}
