package view

import (

)

type Confirmer interface{
	Confirm(info string) bool;
	DisplayAndPause(info string);
}
