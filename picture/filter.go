package picture

import "context"

type SourceFilter struct {
	Provider
	IsAllowed func(*Source) (bool, error)
}

func (p *SourceFilter) GetSourceSet(
	ctx context.Context,
	location string,
) (set []Source, err error) {
	given, err := p.Provider.GetSourceSet(ctx, location)
	if err != nil {
		return nil, err
	}

	for _, source := range given {
		ok, err := p.IsAllowed(&source)
		if err != nil {
			return nil, err
		}
		if ok {
			set = append(set, source)
		}
	}

	return
}
