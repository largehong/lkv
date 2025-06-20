package processor

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
	"github.com/largehong/lkv/command"
	"github.com/sirupsen/logrus"
)

type Processor struct {
	template *template.Template
	dst      string
	after    *command.Command
	src      string
}

func New(tpl *template.Template, src, dst, after string) (processor *Processor) {
	processor = &Processor{
		dst:      dst,
		src:      src,
		template: tpl,
	}
	if after == "" {
		processor.after = nil
	} else {
		processor.after = command.New("sh", "-c", after)
	}
	return processor
}

func (processor *Processor) Redenering() {
	filename := filepath.Join(os.TempDir(), uuid.New().String())
	src, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"template": processor.src,
			"error":    err.Error(),
			"file":     filename,
		}).Error("processor faield to create temp file")
		return
	}
	defer func() {
		src.Close()
		os.Remove(filename)
	}()

	err = processor.template.Execute(src, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"template": processor.src,
			"error":    err.Error(),
		}).Error("failed to render template")
		return
	}

	err = os.Rename(filename, processor.dst)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"template": processor.src,
			"dst":      processor.dst,
			"error":    err.Error(),
			"file":     filename,
		}).Error("failed to mv src to dst")
		return
	}

	if processor.after != nil {
		output, err := processor.after.Run()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"template": processor.src,
				"command":  processor.after.String(),
				"output":   output,
				"err":      err.Error(),
			}).Error("failed to exec after hook")
			return
		}
	}

	logrus.Debugf("rendering %s success", processor.src)
}
