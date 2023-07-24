package backlogcvs

import (
	"context"
	"fmt"
	"hue-n-gnk/csv-to-gitisue/database/dao"
	"hue-n-gnk/csv-to-gitisue/database/entities"
	"hue-n-gnk/csv-to-gitisue/helpers/pointer"
	git "hue-n-gnk/csv-to-gitisue/pkg/github"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const regex = `<.*?>`

func (mc *backlogCsvService) ToGithub() error {
	data, err := download()
	if err != nil {
		return err
	}
	if data == nil {
		return errors.New("input params is missing")
	}
	ctx := context.Background()
	for _, v := range data {
		time.Sleep(3 * time.Second)
		v := v
		if v.AssigneeID != "1149007" && v.AssigneeID != "984140" {
			mc.log.Warn("THE ISSUE NOT BELONG TO US. DO NOTHING!!!")
			continue
		}
		b := mc.db.NulabBacklogIssue
		dt, err := b.Where(b.BacklogID.Eq(v.Id)).First()
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "failed to query data")
			}
		}
		// in case create/update/delete the non exist card
		if errors.Is(err, gorm.ErrRecordNotFound) {
			switch v.StatusId {
			// status resolved
			case "3":
				// do nothing
				mc.log.Warn("THE CARD YOU WANT RESOLVE IS NOT EXIST. DO NOTHING!!!")
				return nil
			case "4":
				// do nothing
				mc.log.Warn("THE CARD YOU WANT CLOSE IS NOT EXIST. DO NOTHING!!!")
				return nil
			default:
				// do create
				err := mc.createGitIssue(ctx, v.Id, v)
				if err != nil {
					return err
				}
			}
			continue
		}
		// in case create/update/delete the exist data
		if pointer.SafeString(dt.GitID) != "" {
			switch v.StatusId {
			case "3", "4":
				err := mc.deleteGitIssue(ctx, dt)
				if err != nil {
					return err
				}
			default:
				err := mc.updateGitIssue(ctx, dt, v)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (mc *backlogCsvService) createGitIssue(ctx context.Context, id string, req CsvStruct) error {
	mc.log.Info("CREATING the ISSUE!!!")

	data := github.IssueRequest{
		Title:  pointer.String(stripHtmlRegex(req.Subject)),
		Body:   stripEOL(req.Description, req.Key),
		Labels: &[]string{"BACKLOG"},
	}
	if req.IssueType != "" {
		*data.Labels = append(*data.Labels, req.IssueType)
	}
	if req.PriorityId != "" {
		*data.Labels = append(*data.Labels, csvToGitPriority(req.PriorityId))
	}
	gt := git.New()

	i, err := gt.CreateIssue(ctx, data)
	if err != nil {
		return err
	}

	ed := &entities.NulabBacklogIssue{
		GitID:     pointer.String(strconv.Itoa(*i)),
		BacklogID: pointer.String(id),
	}

	return mc.db.Transaction(func(tx *dao.Query) error {
		return tx.NulabBacklogIssue.WithContext(ctx).Create(ed)
	})
}

func (mc *backlogCsvService) updateGitIssue(ctx context.Context, e *entities.NulabBacklogIssue, req CsvStruct) error {
	mc.log.Info("UPDATING the ISSUE!!!")

	if pointer.SafeString(e.GitID) == "" {
		return nil
	}
	number, err := strconv.Atoi(pointer.SafeString(e.GitID))
	if err != nil {
		return err
	}

	data := github.IssueRequest{
		Title:  pointer.String(stripHtmlRegex(req.Subject)),
		Body:   stripEOL(req.Description, req.Key),
		Labels: &[]string{"BACKLOG"},
	}
	if req.IssueType != "" {
		*data.Labels = append(*data.Labels, req.IssueType)
	}
	if req.PriorityId != "" {
		*data.Labels = append(*data.Labels, csvToGitPriority(req.PriorityId))
	}
	gt := git.New()

	return gt.UpdateIssue(ctx, number, data)
}

func (mc *backlogCsvService) deleteGitIssue(ctx context.Context, e *entities.NulabBacklogIssue) error {
	mc.log.Info("REMOVING the ISSUE!!!")

	if pointer.SafeString(e.GitID) == "" {
		return nil
	}
	number, err := strconv.Atoi(pointer.SafeString(e.GitID))
	if err != nil {
		return err
	}

	data := github.IssueRequest{
		Labels: &[]string{"BACKLOG", "REMOVED"},
	}
	gt := git.New()

	return gt.UpdateIssue(ctx, number, data)
}

func stripHtmlRegex(s string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}

func stripEOL(s, key string) *string {
	res := strings.ReplaceAll(s, "\\n", "\n")
	res += "\n"
	res += fmt.Sprintf("Ref: https://yhddigital.backlog.com/view/%s", key)
	return &res
}

func csvToGitPriority(s string) string {
	switch s {
	case "1":
		return "Rank S"
	case "2":
		return "Rank A"
	case "3":
		return "Rank B"
	case "4":
		return "Rank:C"
	}
	return ""
}
