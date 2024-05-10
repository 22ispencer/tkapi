package db

type Lab struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func (e *Env) GetLabs() ([]Lab, error) {
	const query = `
	SELECT l.LabID,
		   l.LabName
	FROM [Lab] l
	`
	rows, err := e.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		lab  Lab
		labs = []Lab{}
	)
	for rows.Next() {
		if err := rows.Scan(&lab.Id, &lab.Name); err != nil {
			return nil, err
		}
		labs = append(labs, lab)
	}
	return labs, nil
}
