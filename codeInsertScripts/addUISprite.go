package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"strings"
)

func AppendConst(itemName string) error {
	filePath := "game/interactableUIObjects/constIndex.go"
	lCaser := cases.Lower(language.English)
	lowerCaseInput := lCaser.String(itemName)

	uCaser := cases.Title(language.English)
	newConstName := uCaser.String(lowerCaseInput)

	newConstValue := "\"" + lowerCaseInput + "\""

	// Read the source file
	src, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Create a FileSet
	fset := token.NewFileSet()

	// Parse the source file into an AST
	node, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return err
	}

	// Find the target const declaration
	var targetDecl *ast.GenDecl
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.CONST {
			// You might want to add more specific checks here
			// to ensure you're modifying the correct const block.
			// For example, check for a comment or a specific name.
			targetDecl = genDecl
			break
		}
	}

	if targetDecl == nil {
		fmt.Println("Error: Could not find a const declaration block.")
		return err
	}

	var constType string
	// Attempt to find type from existing const block
	for _, spec := range targetDecl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok && valueSpec.Type != nil {
			if ident, ok2 := valueSpec.Type.(*ast.Ident); ok2 {
				constType = ident.Name
				break
			}
		}
	}

	if constType == "" {
		fmt.Println("Error: Could not determine constant type from existing const block.")
		return err
	}

	// Create the new constant specification
	newSpec := &ast.ValueSpec{
		Names: []*ast.Ident{
			ast.NewIdent(newConstName),
		},
		Type: ast.NewIdent(constType), // Assign the type
		Values: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: newConstValue,
			},
		},
	}

	// Append the new constant specification to the existing const block
	targetDecl.Specs = append(targetDecl.Specs, newSpec)

	// Format the modified AST back into source code
	var buf bytes.Buffer
	err = format.Node(&buf, fset, node)
	if err != nil {
		fmt.Println("Error formatting code:", err)
		return err
	}

	// Write the modified code back to the file
	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	fmt.Printf("Successfully added constant %s to %s\n", newConstName, filePath)
	return nil
}

func AppendJsonWithLocationData(newUIObj string) error {
	var positions = make(map[string]*drawables.SavePositionData)
	spritePosition, err := assets.DataDir.ReadFile("data/spritePosition.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(spritePosition, &positions)
	if err != nil {
		return err
	}

	spData := drawables.SavePositionData{X: 100, Y: 100, Name: newUIObj}

	positions[spData.Name] = &spData

	outputSave, err := json.Marshal(positions)
	if err != nil {
		return err
	}

	err = os.WriteFile("assets/data/spritePosition.json", outputSave, 999)
	if err != nil {
		return err
	}
	return nil
}

func takeInput() (string, error) {
	fmt.Println("Enter Const declaration")
	reader := bufio.NewReader(os.Stdin)
	inputText, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	inputText = strings.TrimSpace(inputText)

	// Sanity check
	if inputText == "" {
		fmt.Println("Empty input, exiting.")
		return "", err
	}

	return inputText, nil
}

func main() {

	inputText, err := takeInput()
	if err != nil {
		log.Fatal(err)
	}

	err = AppendConst(inputText)
	if err != nil {
		log.Fatal(err)
	}

	err = AppendJsonWithLocationData(inputText)
	if err != nil {
		log.Fatal(err)
	}

	mainImgName := "uiSprites/" + inputText + "Main"

	_, err = loaders.LoadImageAssetAsEbitenImage(mainImgName)
	if err != nil {
		log.Printf("no main image found for new ui sprite%q", err)
	} else {
		log.Printf("Image exists for new ui object")
	}

	fmt.Println("successfully added ui sprite constant and appended to sprite location data json ")
}
