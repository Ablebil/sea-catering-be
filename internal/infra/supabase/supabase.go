package supabase

import (
	"fmt"
	"io"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/nedpals/supabase-go"
)

type SupabaseItf interface {
	UploadFile(file io.Reader, bucket string, fileName string, mimeType string) (string, error)
	DeleteFile(bucket string, fileNames []string) error
}

type Supabase struct {
	client *supabase.Client
}

func NewSupabase(conf *conf.Config) SupabaseItf {
	client := supabase.CreateClient(conf.SupabaseURL, conf.SupabaseKey)
	return &Supabase{
		client: client,
	}
}

func safeWrapper(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	f()
	return nil
}

func (s *Supabase) UploadFile(file io.Reader, bucket string, fileName string, mimeType string) (string, error) {
	err := safeWrapper(func() {
		s.client.Storage.From(bucket).Upload(fileName, file, &supabase.FileUploadOptions{
			ContentType: mimeType,
		})
	})
	if err != nil {
		return "", err
	}

	res := s.client.Storage.From(bucket).GetPublicUrl(fileName)
	return res.SignedUrl, nil
}

func (s *Supabase) DeleteFile(bucket string, fileNames []string) error {
	err := safeWrapper(func() {
		s.client.Storage.From(bucket).Remove(fileNames)
	})

	return err
}
