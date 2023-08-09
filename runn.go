package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/k1LoW/runn"
	"github.com/labstack/echo/v4"
)

type RunnResp struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failure int `json:"failure"`
	Skipped int `json:"skipped"`
	Results []struct {
		Id     string `json:"id"`
		Path   string `json:"path"`
		Result string `json:"result"`
		Steps  []any  `json:"steps"`
	} `json:"results"`
}

func handleRunn(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	id, err := saveBook(string(body))
	if err != nil {
		return fmt.Errorf("failed to save book: %w", err)
	}
	c.Logger().Infof("Saved book '%s'", id)

	var opts []runn.Option
	ops, err := runn.Load(bookPath(id), opts...)
	if err != nil {
		return fmt.Errorf("failed to load book: %w", err)
	}

	if err := ops.RunN(context.Background()); err != nil {
		return fmt.Errorf("failed to run book: %w", err)
	}

	if err := saveResult(id, ops.Result().OutJSON); err != nil {
		return fmt.Errorf("failed to save result: %w", err)
	}
	c.Logger().Infof("Saved result for book '%s'", id)

	result, err := os.ReadFile(resultPath(id))
	if err != nil {
		return fmt.Errorf("failed to read result: %w", err)
	}

	var resp RunnResp
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return c.JSON(http.StatusOK, resp)
}

func saveBook(book string) (string, error) {
	id := uuid.New().String()

	f, err := os.Create(bookPath(id))
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.WriteString(book); err != nil {
		return "", err
	}

	return id, nil
}

func saveResult(id string, outJson func(out io.Writer) error) error {
	f, err := os.Create(resultPath(id))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := outJson(f); err != nil {
		return err
	}

	return nil
}

const bookDir = "book/"

func bookPath(id string) string {
	return bookDir + "book-" + id + ".yaml"
}

const resultDir = "result/"

func resultPath(id string) string {
	return resultDir + "result-" + id + ".json"
}
