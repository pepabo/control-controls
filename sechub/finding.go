package sechub

import "fmt"

type FindingGroup struct {
	ControlID string
	Resources FindingResources
}

type FindingGroups []*FindingGroup

func (fgs FindingGroups) ControlIDs() []string {
	ids := []string{}
	for _, fg := range fgs {
		ids = append(ids, fg.ControlID)
	}
	return ids
}

func (fgs FindingGroups) FindByControlID(id string) (*FindingGroup, error) {
	for _, fg := range fgs {
		if fg.ControlID == id {
			return fg, nil
		}
	}
	return nil, fmt.Errorf("not found: %s", id)
}

type FindingResource struct {
	Arn    string
	Status string
	Note   string
}

type FindingResources []*FindingResource

func (frs FindingResources) Arns() []string {
	arns := []string{}
	for _, fr := range frs {
		arns = append(arns, fr.Arn)
	}
	return arns
}

func (frs FindingResources) FindByArn(arn string) (*FindingResource, error) {
	for _, fr := range frs {
		if fr.Arn == arn {
			return fr, nil
		}
	}
	return nil, fmt.Errorf("not found: %s", arn)
}

func intersectFindingGroups(a, b FindingGroups) FindingGroups {
	fgs := FindingGroups{}
	ids := intersect(a.ControlIDs(), b.ControlIDs())
	for _, id := range ids {
		fg := &FindingGroup{ControlID: id}
		afg, _ := a.FindByControlID(id)
		bfg, _ := b.FindByControlID(id)
		if afg == nil || bfg == nil {
			continue
		}
		arns := intersect(afg.Resources.Arns(), bfg.Resources.Arns())
		for _, arn := range arns {
			ar, _ := afg.Resources.FindByArn(arn)
			br, _ := bfg.Resources.FindByArn(arn)
			if ar == nil || br == nil {
				continue
			}
			if ar.Status == br.Status && ar.Note == br.Note {
				fg.Resources = append(fg.Resources, &FindingResource{
					Arn:    arn,
					Status: ar.Status,
					Note:   ar.Note,
				})
			}
		}
		if len(fg.Resources) > 0 {
			fgs = append(fgs, fg)
		}
	}
	return fgs
}

func diffFindingGroups(base, a FindingGroups) FindingGroups {
	fgs := FindingGroups{}
	ids := unique(append(base.ControlIDs(), a.ControlIDs()...))
	for _, id := range ids {
		fg := &FindingGroup{ControlID: id}
		basefg, _ := base.FindByControlID(id)
		afg, _ := a.FindByControlID(id)
		switch {
		case afg == nil:
			// do nothing
		case basefg == nil:
			fg.Resources = afg.Resources
		case basefg != nil && afg != nil:
			arns := unique(append(basefg.Resources.Arns(), afg.Resources.Arns()...))
			for _, arn := range arns {
				baser, _ := basefg.Resources.FindByArn(arn)
				ar, _ := afg.Resources.FindByArn(arn)
				switch {
				case ar == nil:
					// do noting
				case baser == nil:
					fg.Resources = append(fg.Resources, ar)
				case baser != nil && ar != nil:
					if baser.Status != ar.Status || baser.Note != ar.Note {
						fg.Resources = append(fg.Resources, &FindingResource{
							Arn:    arn,
							Status: ar.Status,
							Note:   ar.Note,
						})
					}
				}
			}
		}
		if len(fg.Resources) > 0 {
			fgs = append(fgs, fg)
		}
	}
	return fgs
}

func overlayFindingGroups(base, overlay FindingGroups) FindingGroups {
	for _, ofg := range overlay {
		basefg, _ := base.FindByControlID(ofg.ControlID)
		if basefg == nil {
			base = append(base, ofg)
			continue
		}
		for _, r := range ofg.Resources {
			baser, _ := basefg.Resources.FindByArn(r.Arn)
			if baser == nil {
				basefg.Resources = append(basefg.Resources, r)
				continue
			}
			baser.Status = r.Status
			baser.Note = r.Note
		}
	}
	return base
}
