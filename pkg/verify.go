package pkg

var (
	NameVerify               = Rules{"Name": {NotEmpty(), RegexpMatch("^(男人|女人)?$")}}
)
