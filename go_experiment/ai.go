package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const SETTINGS_FILE_PATH = "settings.json"

type Api struct {
	api_key string `json:"api_key"'`
}

func retrieveSettings(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("%s", err)
		return ""
	}
	defer f.Close()	

	/* Read JSON */
	bytes_data, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("%s", err)
		return ""
	}

	var api Api
	if err := json.Unmarshal(bytes_data, &api); err != nil {
		fmt.Printf("%s", err)
		return ""
	}
	
	return api.api_key
}

func googleModel(pdf_text string, material_type bool, answer_mode bool) string {
	/* Setup Model */
	api_key := retrieveSettings(SETTINGS_FILE_PATH)
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	model := client.GenerativeModel("gemini-1.5-flash")
	
	/* Setup Prompt */
	switch {
		case material_type && !answer_mode: /* Practice exam without answers */
			prompt := fmt.Sprintf("Below is a practice exam. I want you to create another practice exam that is similar to this but harder. Here is the practice exam:\n%s", pdf_text)
			/* Gen Practice exam without answers */
			response, err := model.GenerateContent(ctx, genai.Text(prompt))
			return response
		case material_type && answer_mode: /* Generate Answers */
			prompt := fmt.Sprintf("Below is a practice exam. I want you to create another practice exam that is similar to this but harder. Here is the practice exam:\n%s", pdf_text)
			response, err := model.GenerateContent(ctx, genai.Text(prompt))
			return response
		case !material_type:
			prompt := fmt.Sprintf("Below are lecture notes. I want you to generate me markdown text that summarises these lecture notes but don't make me miss out on any important information. You can include your own knowledge in this markdown text. Make these summaries simple so that someone who has no knowledge can understand this topic e.g. using the feynman technique. Here are the lecture notes:\n%s", pdf_text)
			response, err := model.GenerateContent(ctx, genai.Text(prompt))
			return response
	}
	
	return ""
}




