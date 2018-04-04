package coordinats

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"github.com/pebbe/go-proj-4/proj"
)

type MSK struct {
	x, y float64
}

func CreateMSKCoordinate(x, y float64) *MSK {
	var msk MSK
	msk.x = x
	msk.y = y
	return &msk
}

func ParseMSKCoordinate(str string) (*MSK, error) {
	re, err := regexp.Compile("(\\d+[\\.,]{0,1}\\d+)[\\D]+(\\d+[\\.,]{0,1}\\d+)")
	if err != nil {
		return nil, err
	}
	coord := re.FindAllString(str, -1)
	if len(coord) == 0 || len(coord) > 1 {
		return nil, errors.New("Parse error: too many matches")
	}
	mskCoordXY := re.FindStringSubmatch(str)
	mskXStr := mskCoordXY[1]
	mskYStr := mskCoordXY[2]
	fmt.Println(mskXStr, mskYStr)
	mskXStr = strings.Replace(mskXStr, ",", ".", -1)
	mskYStr = strings.Replace(mskYStr, ",", ".", -1)
	mskXStr = strings.TrimSpace(mskXStr)
	mskYStr = strings.TrimSpace(mskYStr)
	mskX, err := strconv.ParseFloat(mskXStr, 64)
	mskY, err := strconv.ParseFloat(mskYStr, 64)
	if err != nil {
		return nil, errors.New("Parse coordinate error: wrong coordinate values")
	}
	mskStruct := CreateMSKCoordinate(mskX, mskY)
	return mskStruct, nil
}

func (msk *MSK) MSKToWGS() (*WGS, error) {
	wgsParams, err := proj.NewProj("+proj=longlat +datum=WGS84 +ellps=WGS84 +no_defs +axis=neu")
	defer wgsParams.Close()
	if err != nil {
		return nil, err
	}
	mskParams, err := proj.NewProj("+proj=tmerc +ellps=krass +towgs84=23.57,-140.95,-79.8,0,0.35,0.79,-0.22 +units=m +lon_0=21.45 +lat_0=0 +k_0=1 +x_0=1250000 +y_0=-5711057.628 +no_defs")
	defer mskParams.Close()
	if err != nil {
		return nil, err
	}
	lat, long, err := proj.Transform2(mskParams, wgsParams, msk.y, msk.x)
	if err != nil {
		return nil, err
	}
	var wgsStruct WGS
	wgsStruct.Lat = proj.RadToDeg(lat)
	wgsStruct.Long = proj.RadToDeg(long)
	return &wgsStruct, nil
}
