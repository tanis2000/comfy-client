package objectinfo

import (
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "io"
    "os"
    "testing"
)

func TestUnmarshalOk(t *testing.T) {
    f, err := os.Open("./testdata/data.json")
    assert.NoError(t, err)
    buf, err := io.ReadAll(f)
    assert.NoError(t, err)
    info := &ObjectInfo{}
    err = json.Unmarshal(buf, &info)
    assert.NoError(t, err)
}

func TestGetByTypeOk(t *testing.T) {
    f, err := os.Open("./testdata/data.json")
    assert.NoError(t, err)
    buf, err := io.ReadAll(f)
    assert.NoError(t, err)
    info := &ObjectInfo{}
    err = json.Unmarshal(buf, &info)
    assert.NoError(t, err)
    ckpt := info.GetByType("CheckpointLoaderSimple")
    assert.NotNil(t, ckpt)
}

func TestGetByTypeNotExists(t *testing.T) {
    f, err := os.Open("./testdata/data.json")
    assert.NoError(t, err)
    buf, err := io.ReadAll(f)
    assert.NoError(t, err)
    info := &ObjectInfo{}
    err = json.Unmarshal(buf, &info)
    assert.NoError(t, err)
    ckpt := info.GetByType("BadNameForCheckpointType")
    assert.Nil(t, ckpt)
}

func TestGetListOfCheckpointsOk(t *testing.T) {
    f, err := os.Open("./testdata/data.json")
    assert.NoError(t, err)
    buf, err := io.ReadAll(f)
    assert.NoError(t, err)
    info := &ObjectInfo{}
    err = json.Unmarshal(buf, &info)
    assert.NoError(t, err)
    ckpts, err := info.GetListOfCheckpoints()
    assert.NoError(t, err)
    assert.NotNil(t, ckpts)
    assert.Len(t, ckpts, 51)
}

func TestGetListOfImagesOk(t *testing.T) {
    f, err := os.Open("./testdata/data.json")
    assert.NoError(t, err)
    buf, err := io.ReadAll(f)
    assert.NoError(t, err)
    info := &ObjectInfo{}
    err = json.Unmarshal(buf, &info)
    assert.NoError(t, err)
    ckpts, err := info.GetListOfImages()
    assert.NoError(t, err)
    assert.NotNil(t, ckpts)
    assert.Len(t, ckpts, 15)
}
