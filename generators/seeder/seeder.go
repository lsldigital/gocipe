package seeder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fluxynet/gocipe/util"
)

// Generate returns generated database schema creation code
func Generate(r *util.Recipe) {

	data := util.GenerataSeeds(r)
	rankings, _ := json.Marshal(data)
	err := ioutil.WriteFile(util.WorkingDir+"output.json", rankings, 0644)

	fmt.Printf("%+v", rankings)
	fmt.Printf("%+v", err)

}
