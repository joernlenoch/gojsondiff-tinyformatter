package formatter

import (
  "fmt"
  "errors"

  diff "github.com/yudai/gojsondiff"
)

const (
  DeltaDelete   = 0
  DeltaTextDiff = 2
  DeltaMove     = 3
)

func NewTinyFormatter() *TinyFormatter {
  return &TinyFormatter{}
}

type TinyFormatter struct {}

func (f *TinyFormatter) FormatObject(deltas []diff.Delta) (deltaJson map[string]interface{}, err error) {
  deltaJson = map[string]interface{}{}
  for _, delta := range deltas {
    switch delta.(type) {
    case *diff.Object:
      d := delta.(*diff.Object)
      deltaJson[d.Position.String()], err = f.FormatObject(d.Deltas)
      if err != nil {
        return nil, err
      }
    case *diff.Array:
      d := delta.(*diff.Array)
      deltaJson[d.Position.String()], err = f.FormatArray(d.Deltas)
      if err != nil {
        return nil, err
      }
    case *diff.Added:
      d := delta.(*diff.Added)
      deltaJson[d.PostPosition().String()] = []interface{}{d.Value}
    case *diff.Modified:
      d := delta.(*diff.Modified)
      deltaJson[d.PostPosition().String()] = []interface{}{0, d.NewValue}
    case *diff.TextDiff:
      d := delta.(*diff.TextDiff)
      deltaJson[d.PostPosition().String()] = []interface{}{d.DiffString(), 0, DeltaTextDiff}
    case *diff.Deleted:
      d := delta.(*diff.Deleted)
      deltaJson[d.PrePosition().String()] = []interface{}{0, 0, DeltaDelete}
    case *diff.Moved:
      return nil, errors.New("delta type 'Move' is not supported in objects")
    default:
      return nil, errors.New(fmt.Sprintf("unknown delta type detected: %#v", delta))
    }
  }
  return
}

func (f *TinyFormatter) FormatArray(deltas []diff.Delta) (deltaJson map[string]interface{}, err error) {
  deltaJson = map[string]interface{}{
    "_t": "a",
  }
  for _, delta := range deltas {
    switch d := delta.(type) {
    case *diff.Object:
      deltaJson[d.Position.String()], err = f.FormatObject(d.Deltas)
      if err != nil {
        return nil, err
      }
    case *diff.Array:
      deltaJson[d.Position.String()], err = f.FormatArray(d.Deltas)
      if err != nil {
        return nil, err
      }
    case *diff.Added:
      deltaJson[d.PostPosition().String()] = []interface{}{d.Value}
    case *diff.Modified:
      deltaJson[d.PostPosition().String()] = []interface{}{0, d.NewValue}
    case *diff.TextDiff:
      deltaJson[d.PostPosition().String()] = []interface{}{d.DiffString(), 0, DeltaTextDiff}
    case *diff.Deleted:
      deltaJson["_"+d.PrePosition().String()] = []interface{}{0, 0, DeltaDelete}
    case *diff.Moved:
      deltaJson["_"+d.PrePosition().String()] = []interface{}{d.Value, d.PostPosition(), DeltaMove}
    default:
      return nil, errors.New(fmt.Sprintf("unknown delta type detected: %#v", delta))
    }
  }
  return
}

