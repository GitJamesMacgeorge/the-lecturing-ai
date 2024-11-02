import os, json
import ai as ai
from pypdf import PdfReader
from markdown_pdf import MarkdownPdf, Section

# CONSTANTS
ORIGINAL_PRACTICE_EXAMS_DIR = "originalexams"
ORIGINAL_NOTES_DR = "originalnotes"
GENERATED_EXAMS_DIR = "generatedexams"
GENERATED_NOTES_DIR = "generatednotes"

def extractPDFtext(file_path:str, material_type:bool):
    if material_type:
        file_path = f"{ORIGINAL_PRACTICE_EXAMS_DIR}/{file_path}.pdf"
    else:
        file_path = f"{ORIGINAL_NOTES_DR}/{file_path}.pdf"

    if not os.path.isfile(file_path):
        raise Exception(f"Error: {file_path} doesn't exist")

    pdf_text = f""
    reader = PdfReader(file_path)
    for page in reader.pages:
        page_text = page.extract_text()
        page_text = page_text.replace("\t", " ")
        pdf_text += page_text
    
    pdf_text = pdf_text.replace("FINC2012", "subject").replace("finc2012", "subject").replace("Unviversity of Sydney", "").replace("USYD", "")
    return pdf_text

def pdfConversion(markdown_content:str, pdf_name:str, material_type:str):
    save_path = f"{GENERATED_NOTES_DIR}/{pdf_name}.pdf"
    if material_type:
        save_path = f"{GENERATED_EXAMS_DIR}/{pdf_name}.pdf"
    
    # Generate PDF 
    pdf = MarkdownPdf()
    pdf.add_section(Section(markdown_content, toc=False))
    pdf.save(save_path)

    # Check if file now exists
    if os.path.isfile(save_path):
        print(f"{pdf_name}.pdf has been generated")

def handleUserResponse():
    user_input = input()
    user_input = user_input.strip("\n")
    user_input = user_input.strip()
    return user_input

def promptUser():
    user_response = {
        "filename": None,
        "material_type": None,
        "answer_mode": None,
        "desired_name": None
    }

    print("Enter the PDF filename: ", end="")
    user_response["filename"] = handleUserResponse()
    print("\nIs this a Practice Exam? (Y/n) ", end="")
    material_type = handleUserResponse().lower().strip("\n")
    if material_type == "y":
        user_response["material_type"] = True
        # Ask for answers
        print("\nDo you want the answers to this material? (Y/n) ", end="")
        user_input = handleUserResponse().lower().strip("\n").strip()
        if user_input == "y":
            user_response["answer_mode"] = True
        elif user_input == "n":
            user_response["answer_mode"] = False 
        else:
            print(user_input)
            raise Exception("Error: Invalid user input for answer mode")
    elif material_type == "n":
        user_response["material_type"] = False
    else:
        raise Exception("Error: Invalid user input for material type")

    # Handle Desired Name
    print(f"\nDesired filename: ", end="")
    user_response["desired_name"] = handleUserResponse()
    
    return user_response

def initAi(pdf_text:str, material_type:bool, answer_mode:bool):
    # Generate Response from model
    model_response = ai.googleGem(pdf_text, material_type, answer_mode, 1)
    return model_response

def main():
    # Collect User Response
    user_response = promptUser()

    # Retrieve PDF text
    pdf_text = extractPDFtext(user_response["filename"], user_response["material_type"])

    # Activate LLM
    if not user_response["answer_mode"]:
        model_response = initAi(pdf_text, user_response["material_type"], user_response["answer_mode"]) 
        pdfConversion(model_response, user_response["desired_name"], user_response["material_type"])
    else:
        exam_response = initAi(pdf_text, user_response["material_type"], False) 
        answer_response = initAi(exam_response, user_response["material_type"], True)
        pdfConversion(exam_response, user_response["desired_name"], user_response["material_type"])
        pdfConversion(answer_response, user_response["desired_name"] + " (ANSWERS)", user_response["material_type"])

if __name__ == "__main__":
    main()
