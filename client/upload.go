package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "strconv"
)

type ImageType string

const (
    InputImageType  ImageType = "input"
    TempImageType   ImageType = "temp"
    OutputImageType ImageType = "output"
)

func (c *Client) Upload(r io.Reader, filename string, overwrite bool, fileType ImageType, subfolder string) (string, error) {
    var reqBody bytes.Buffer
    writer := multipart.NewWriter(&reqBody)
    formFile, err := writer.CreateFormFile("image", filename)
    if err != nil {
        return "", err
    }
    _, err = io.Copy(formFile, r)
    if err != nil {
        return "", err
    }
    err = writer.WriteField("overwrite", strconv.FormatBool(overwrite))
    if err != nil {
        return "", err
    }
    err = writer.WriteField("type", string(fileType))
    if err != nil {
        return "", err
    }
    if subfolder != "" {
        err = writer.WriteField("subfolder", subfolder)
        if err != nil {
            return "", err
        }
    }
    err = writer.Close()
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", c.buildUrl("/upload/image"), &reqBody)
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("upload failed with status code %d", resp.StatusCode)
    }

    var data map[string]interface{}
    if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return "", err
    }
    name, ok := data["name"].(string)
    if !ok {
        return "", fmt.Errorf("invalid response from upload")
    }
    return name, nil
}
