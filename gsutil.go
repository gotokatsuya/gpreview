package gpreview

import (
	"fmt"
	"os"
	"os/exec"
)

func CmdRun(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ShowReviews() error {
	if err := CmdRun("gsutil", "ls", fmt.Sprintf("gs://%s/reviews/", GPReview.BucketID)); err != nil {
		return err
	}
	return nil
}

func CopyReviews() error {
	if err := CmdRun("gsutil", "cp", "-r", fmt.Sprintf("gs://%s/reviews/reviews_%s_*", GPReview.BucketID, GPReview.PackageName)); err != nil {
		return err
	}
	return nil
}
