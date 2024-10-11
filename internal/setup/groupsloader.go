package setup

import (
	route_Test1 "github.com/krittakondev/goapisuit/internal/routes/test1"
	route_Test1Test2 "github.com/krittakondev/goapisuit/internal/routes/test1/test2"
	"github.com/krittakondev/goapisuit"
)

func GroupsSetup(suit *goapisuit.Suit){
	suit.SetupGroups(suit.Config.ApiPrefix+"/test1", &route_Test1.Route{})
	suit.SetupGroups(suit.Config.ApiPrefix+"/test1/test2", &route_Test1Test2.Route{})
}
