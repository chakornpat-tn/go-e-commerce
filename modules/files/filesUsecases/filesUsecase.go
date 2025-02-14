package filesUsecases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/files"
)

type IFilesUsecase interface {
	UpLoadToGCP(req []*files.FileReq) ([]*files.FileRes, error)
	DeleteFileOnGCP(req []*files.DeleteFileReq) error
}

type filesUsecase struct {
	cfg config.IConfig
}

func NewFilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

func (u *filesUsecase) UpLoadToGCP(req []*files.FileReq) ([]*files.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	jobCh := make(chan *files.FileReq, len(req))
	resCh := make(chan *files.FileRes, len(req))
	errCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0)

	for _, r := range req {
		jobCh <- r
	}
	close(jobCh)

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go u.upLoadWorker(ctx, client, jobCh, resCh, errCh)
	}

	for l := 0; l < len(req); l++ {
		err := <-errCh
		if err != nil {
			return nil, err
		}
		result := <-resCh
		res = append(res, result)
	}

	return res, nil
}

func (u *filesUsecase) upLoadWorker(ctx context.Context, client *storage.Client, jobs <-chan *files.FileReq, result chan<- *files.FileRes, errs chan<- error) {
	for job := range jobs {

		container, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}

		b, err := ioutil.ReadAll(container)
		if err != nil {
			errs <- err
			return
		}

		buf := bytes.NewBuffer(b)

		// Upload an object with storage.Writer.
		wc := client.Bucket(u.cfg.APP().GCPBucket()).Object(job.Destination).NewWriter(ctx)

		if _, err = io.Copy(wc, buf); err != nil {
			errs <- fmt.Errorf("io.Copy: %w", err)
			return
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errs <- fmt.Errorf("Writer.Close: %w", err)
			return
		}
		fmt.Printf("%v uploaded to %v.\n", job.FileName, job.Destination)

		newFile := &filePub{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      "https://storage.googleapis.com/" + u.cfg.APP().GCPBucket() + "/" + job.Destination,
			},
			bucket:      u.cfg.APP().GCPBucket(),
			destination: job.Destination,
		}

		if err := newFile.setPublic(ctx, client); err != nil {
			errs <- err
			return
		}

		errs <- nil
		result <- newFile.file
	}
}

type filePub struct {
	bucket      string
	destination string
	file        *files.FileRes
}

func (f *filePub) setPublic(ctx context.Context, client *storage.Client) error {
	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %w", err)
	}
	fmt.Printf("Bucket %v is now publicly accessible.\n", f.destination)
	return nil

}

func (u *filesUsecase) deleteFile(ctx context.Context, client *storage.Client, jobs <-chan *files.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		o := client.Bucket(u.cfg.APP().GCPBucket()).Object(job.Destination)

		attrs, err := o.Attrs(ctx)
		if err != nil {
			errs <- fmt.Errorf("object.Attrs: %w", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errs <- fmt.Errorf("Object(%q).Delete: %w", job.Destination, err)
			return
		}
		fmt.Printf("Blob %v deleted.\n", job.Destination)

		errs <- nil
	}

}

func (u *filesUsecase) DeleteFileOnGCP(req []*files.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	jobCh := make(chan *files.DeleteFileReq, len(req))
	errCh := make(chan error, len(req))

	for _, r := range req {
		jobCh <- r
	}
	close(jobCh)

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go u.deleteFile(ctx, client, jobCh, errCh)
	}

	for l := 0; l < len(req); l++ {
		err := <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}
