package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dslipak/pdf"
	fpdf "github.com/go-pdf/fpdf"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
)

const ORIGINAL_PRACTICE_EXAMS_DIR = "originalexams"
const ORIGINAL_NOTES_DIR = "originalnotes"
const GENERATED_EXAMS_DIR = "generatedexams"
const GENERATED_NOTES_DIR = "generatednotes"

type Userrequest struct {
	filename string
	material_type bool
	answer_mode bool 
	desired_name string
}

func pdfConversion(pdf_text string, desired_name string, material_type bool) {
	save_path := fmt.Sprintf("%s/%s.pdf", GENERATED_NOTES_DIR, desired_name)
	if material_type {
		save_path = fmt.Sprintf("%s/%s.pdf", GENERATED_EXAMS_DIR, desired_name)
	}

	/* Converts MD to HTML */
	unsafe_data := blackfriday.Run([]byte(pdf_text))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe_data)
	
	/* Generates PDF file */
	pdf_file := fpdf.New("P", "mm", "A4", "")
	pdf_file.AddPage()
	pdf_file.SetFont("Arial", "", 11)
	pdf_file.OutputFileAndClose(save_path)

}

func doesFileExist(filename string) bool {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}
	return true
}

func readPdf(path string) string {
	r, err := pdf.Open(path)
	// remember close file
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
    b, err := r.GetPlainText()
    if err != nil {
        return ""
    }
    buf.ReadFrom(b)
	return buf.String()
}

func extractPDFtext(filepath string, material_type bool) string {
	if material_type {
		filepath = fmt.Sprintf("%s/%s.pdf", ORIGINAL_PRACTICE_EXAMS_DIR, filepath)
	} else {
		filepath = fmt.Sprintf("%s/%s.pdf", ORIGINAL_NOTES_DIR, filepath)
	}
	
	if !doesFileExist(filepath) {
		errors.New("Error: Output pdf file already exist")
		fmt.Printf("%s doesnt exist", filepath)
		return ""
	}

	/* Retrieve PDF text */
	x := readPdf(filepath)
	fmt.Printf("%s", x)
	return ""
}

func getUserReply(prompt string) string {
	var user_reply string
	fmt.Printf("%s", prompt)
	fmt.Scan(&user_reply)
	return user_reply
}

func promptUser() Userrequest {
	var filename, material_type, answer_mode, desired_name string 
	filename = getUserReply("Enter the PDF filename: ")
	material_type = strings.ToLower(strings.TrimSpace(getUserReply("Is this a Practice Exam? (Y/n)")))
	if material_type == "y" {
		answer_mode = strings.ToLower(strings.TrimSpace(getUserReply("Do you want answers for this? (Y/n)")))
	}
	desired_name = strings.ToLower(strings.TrimSpace(getUserReply("What do you want to call this pdf? ")))

	/* Generate Struct */
	actual_type := false 
	if material_type == "y" {
		actual_type = true
	}
	
	actual_mode := false 
	if answer_mode == "y" {
		actual_mode = true
	}
	user_request := Userrequest {
		filename: filename,
		material_type: actual_type,
		answer_mode: actual_mode,
		desired_name: desired_name,
	}

	return user_request
}

func main() {
	user_request := promptUser()
	
	/* Extract text from PDF */
    pdf_text := extractPDFtext(user_request.filename, user_request.material_type)

	/* Parse text to AI model */
	model_response := ""
	if !user_request.material_type {
		/* User wants generated lecture notes */
		model_response = googleModel(pdf_text, user_request.material_type, user_request.answer_mode)
	}
}
