package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

// ---------------- Types ----------------

type DMS struct {
	Degrees int
	Minutes int
	Seconds int
}

func (d DMS) String() string {
	return fmt.Sprintf("%d° %d' %d\"", d.Degrees, d.Minutes, d.Seconds)
}

type VargaResult struct {
	Varga string `json:"Varga"`
	VargaID int `json:VargaID`
	SignIndex int `json:SignIndex`
	Sign  string `json:"Sign"`
	DMS   DMS    `json:"DMS"`
}

// Pie defines the angular boundaries of a slice
type Rashi struct {
	Name  string
	Start float64
	End   float64
}

// ---------------- Globals ----------------

//{Rashiname, start degrees - inclusive, end degrees - exclusive}          >= start and < end
var Rashis = []Rashi{
	{"Mesha",0,30},
	{"Vrushabha",30,60},
	{"Mithuna",60,90},
	{"Karka",90,120},
	{"Simha",120,150},
	{"Kanya",150,180},
	{"Tula",180,210},
	{"Vrischika", 210, 240},
	{"Dhanusha",240,270},
	{"Makara",270,300},
	{"Kumbha",300,330},
	{"Meen",330,360},
}


var vargas  []VargaResult

var signs = []string{
	"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
	"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

// ---------------- Helpers ----------------

func toDMS(deg float64) DMS {
	d := int(deg)
	m := int((deg - float64(d)) * 60)
	s := int(((deg - float64(d)) * 60 - float64(m)) * 60)
	return DMS{Degrees: d, Minutes: m, Seconds: s}
}

// ---------------- Divisional Functions ----------------

// D1: Rasi
func D1(baseSignIndex int, posInSign float64) VargaResult {
	return VargaResult{"D1 Rasi", 1, baseSignIndex, signs[baseSignIndex], toDMS(posInSign)}
}


//D2 Hora
func D2HD(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0/2
	amsaIndex := int(posInSign / amsaSize)

	var mapped int	
	//since 0 based, this is actually an odd rasi
	if (baseSignIndex%2) == 0 {
		if amsaIndex == 1 {
			mapped = 3 //Moon's Hora - Cancer		
		} else {
			mapped = 4 //Sun's Hora - Leo
		}
	} else {

		if amsaIndex == 1 {
			mapped = 4 //Sun's Hora - Cancer		
		} else {
			mapped = 3 //Moon's Hora - Leo
		}
	}

	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 27
	fmt.Printf("D2HD Hora baseSignIndex: %d , posInSign: %f amsaSize: %f amsaIndex: %d mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, amsaIndex, mapped, degInAmsa)
	return VargaResult{"D2HD Hora", 2, mapped, signs[mapped], toDMS(degInAmsa)}

}

// D3: Drekkana
//verified with some tests
func D3(baseSignIndex int, posInSign float64) VargaResult {
//	amsaSize := 30.0 / 3
	amsaSize := 10.0
	amsaIndex := int(posInSign / amsaSize)
	offsets := []int{0, 4, 8}
	mapped := (baseSignIndex + offsets[amsaIndex]) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 3
//	fmt.Println("baseSignIndex: ", baseSignIndex, " posInSign: ", posInSign,  " amsaIndex:", amsaIndex, " mapped: ", mapped, " degInAmsa: ", degInAmsa)
	return VargaResult{"D3 Drekkana", 3,  mapped, signs[mapped], toDMS(degInAmsa)}
}

//D4 HD as per PVRN
func D4HD(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 7.5 
	amsaIndex := int(posInSign / amsaSize)
	offsets := []int{0, 3, 6, 9}
	mapped := (baseSignIndex + offsets[amsaIndex]) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 3
//	fmt.Println("baseSignIndex: ", baseSignIndex, " posInSign: ", posInSign,  " amsaIndex:", amsaIndex, " mapped: ", mapped, " degInAmsa: ", degInAmsa)
	return VargaResult{"D4 HD Chaturthamsa", 4, mapped, signs[mapped], toDMS(degInAmsa)}
}

//D5 as per PVRN
//Odd rasis amsas go in Aries, Aquarius, Saggitarius, Gemini and Libra in that order.
//Even rasis amsas go in Tauris, Virgo, Pisces, Capricon and Scorpio in that order

func D5HD(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 6.0 
	amsaIndex := int(posInSign / amsaSize)

	var mapped int
	//since 0 based, this is actually an odd rasi
	//if math.Mod(baseSignIndex,2) == 0 {
	if (baseSignIndex%2) == 0 {
		switch amsaIndex {
			case 0:
				mapped = 0 //Aries
			case 1:
				mapped = 10 //Aquarius
			case 2:
				mapped = 8 // Saggitarius
			case 3:
				mapped = 2 // Gemini
			case 4:
				mapped = 6 //Libra
		}
	} else { //even rasis
		switch amsaIndex {
			case 0:
				mapped = 1 // Taurus
			case 1:
				mapped = 5 //Virgo
			case 2:
				mapped = 11 // Pisces
			case 3:
				mapped = 9 // Capricorn
			case 4:
				mapped = 7 //Scorpio
		}
	}
	
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 6
//	fmt.Println("baseSignIndex: ", baseSignIndex, " posInSign: ", posInSign,  " amsaIndex:", amsaIndex, " mapped: ", mapped, " degInAmsa: ", degInAmsa)
	return VargaResult{"D5 HD Panchamsa", 5, mapped, signs[mapped], toDMS(degInAmsa)}
}


//D6 as per PVRN
//start from Aries for odd rasis and Libra for even
func D6HD(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 5.0 
	amsaIndex := int(posInSign / amsaSize)

	var mapped int
	//since 0 based, this is actually an odd rasi
	//if math.Mod(baseSignIndex,2) == 0 {
	if (baseSignIndex%2) == 0 {
		mapped = amsaIndex  // Aries is 0, so nothing to add  
		
	} else { //even rasis
			
		mapped = amsaIndex + 6 // Libra starting point for even rasis
	}
	
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 5
//	fmt.Println("baseSignIndex: ", baseSignIndex, " posInSign: ", posInSign,  " amsaIndex:", amsaIndex, " mapped: ", mapped, " degInAmsa: ", degInAmsa)
	return VargaResult{"D6 HD Shasthamsa", 6, mapped, signs[mapped], toDMS(degInAmsa)}
}


// D7: Saptamsa
func D7(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 7
	amsaIndex := int(posInSign / amsaSize)
	start := baseSignIndex
	if baseSignIndex%2 != 0 { // even signs
		start = (baseSignIndex + 6) % 12
	}
	mapped := (start + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 7

	fmt.Printf(" D7 Saptamsa baseSignIndex: %d , posInSign: %f amsaSize: %f mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, mapped, degInAmsa)
	return VargaResult{"D7 Saptamsa", 7, mapped, signs[mapped], toDMS(degInAmsa)}
}

// D7 HD: Saptamsa
//FIXME - logic validation needed for all cases. example - 210 degrees breaks  
func D7HD(baseSignIndex int, posInSign float64) VargaResult {
//	const DIVSIZE int = 7 //this approach causes some issues

	amsaSize := 30.0 / 7 
	amsaIndex := int(posInSign / amsaSize)

	//start := baseSignIndex
	var mapped int
	if baseSignIndex%2 != 0 { // even signs because baseSignIndex is 0 based
		mapped = baseSignIndex + amsaIndex + 6
	} else {
		mapped = baseSignIndex + amsaIndex
	}
	if mapped > 11 {
		mapped = mapped - 11
	} 
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 7 


	fmt.Printf(" D7 HD Saptamsa baseSignIndex: %d , posInSign: %f amsaSize: %f mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, mapped, degInAmsa)

	return VargaResult{"D7 HD Saptamsa", 7, mapped, signs[mapped], toDMS(degInAmsa)}
}


// D8 HD: Saptamsa
// tested with 3 cases
func D8HD(baseSignIndex int, posInSign float64) VargaResult {

	amsaSize := 30.0 / 8 
	amsaIndex := int(posInSign / amsaSize)

	//start := baseSignIndex
	var mapped int

	switch baseSignIndex { //0 based
	case 0,3,6,9: //movable
		mapped = 0 + amsaIndex //start from Aries
	case 1,4,7,10: //fixed
		mapped = 4 + amsaIndex //start from Leo
	case 2,5,8,11: //dual
		mapped = 8 + amsaIndex
	}

	mapped = mapped % 12
	fmt.Printf("baseSignIndex: %d , mapped: %d \n",baseSignIndex, mapped)
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 8 
	return VargaResult{"D8 HD Ashtamsa", 8, mapped, signs[mapped], toDMS(degInAmsa)}
}

// D9: Navamsa
func D9(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 9  //3.20

	amsaIndex := int(posInSign / amsaSize)

	var start int
	switch baseSignIndex {
	case 0, 3, 6, 9: // movable
		start = baseSignIndex
	case 1, 4, 7, 10: // fixed
		start = (baseSignIndex + 8) % 12
	default: // dual
		start = (baseSignIndex + 4) % 12
	}
	mapped := (start + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 9
	return VargaResult{"D9 Navamsa", 9, mapped, signs[mapped], toDMS(degInAmsa)}
}


//D9 - my logic based on PVRN 

func D9HD(baseSignIndex int, posInD1 float64) VargaResult {
	//fmt.Printf("D9HD baseSignIndex: %d\n", baseSignIndex)

	amsaSize := 30.0 / 9  //3.33 using degree decimal notation.

	amsaIndex := int(posInD1 / amsaSize)

	//baseSignIndex is 0 based
	var start int
	switch baseSignIndex {
	case 0, 4, 8: // fiery 
		start = 0 
	case 1, 5, 9: // earthy 
		start = 9
	case 2, 6, 10: // windy
		start = 6
	case 3, 7, 11: //watery
		start = 3 
	default: // problem case 
		start = 100
	}
	//sign in Navamsa
	mapped := (start + amsaIndex) % 12
	//degrees in amsa
	degInAmsa := (posInD1 - float64(amsaIndex)*amsaSize) * 9
	//degInAmsa :=  (posInD1 - (float64(amsaIndex - 1) * amsaSize)) * 9  

	//fmt.Printf("amsaSize : %f amsaIndex : %d  degInAmsa :  %f \n", amsaSize, amsaIndex, degInAmsa)

	return VargaResult{"D9 HD", 9, mapped, signs[mapped], toDMS(degInAmsa)}
}



// D9 (Elemental variant): 
// Fiery signs start from Aries, Earthy from Capricorn, Airy from Libra, Watery from Cancer
//same logic as above from pvrn - this function also works correctly
/*
func D9_Parasara(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 9
	amsaIndex := int(posInSign / amsaSize)

	// Determine element
	var start int
	switch baseSignIndex {
	case 0, 4, 8: // Aries, Leo, Sagittarius -> Fire
		start = 0 // Aries
	case 1, 5, 9: // Taurus, Virgo, Capricorn -> Earth
		start = 9 // Capricorn
	case 2, 6, 10: // Gemini, Libra, Aquarius -> Air
		start = 6 // Libra
	case 3, 7, 11: // Cancer, Scorpio, Pisces -> Water
		start = 3 // Cancer
	}

	mapped := (start + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 9
	return VargaResult{"D9 Navamsa (Parasara)", signs[mapped], toDMS(degInAmsa)}
}
*/


// D10: Dasamsa
//tested 3 cases
func D10(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 10
	amsaIndex := int(posInSign / amsaSize)
	start := baseSignIndex
	if baseSignIndex%2 != 0 { // even
		start = (baseSignIndex + 8) % 12
	}
	mapped := (start + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 10
	return VargaResult{"D10 Dasamsa", 10, mapped, signs[mapped], toDMS(degInAmsa)}
}

// D11: Rudramsa / Ekadasama
//count baseSignIndex from Aries, go in opposite direction from Aries for equal number.
//this is starting point. count number of amsas from this.
func D11HD(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 11

	amsaIndex := int(posInSign / amsaSize)

	//fmt.Printf("D11HD amsaSize: %s, AmsaIndex: %d  \n", toDMS(amsaSize), amsaIndex)
	start :=  12 - baseSignIndex

	var mapped int
	mapped = (start + amsaIndex)
	if (mapped > 11 ) {
		mapped =  mapped - 12
	}

	//fmt.Printf("degInAmsa %f \n", posInSign -  (   ))
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 11

	fmt.Printf(" D11HD Rudramsa/Ekadasama  baseSignIndex: %d , posInSign: %f amsaSize: %f mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, mapped, degInAmsa)

	return VargaResult{"D11 HD Rudramsa", 11, mapped, signs[mapped], toDMS(degInAmsa)}

//	return VargaResult{"D11 HD Rudramsa", 0, signs[0], toDMS(2.5)}

}


// D12: Dvadasamsa
func D12(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 12
	amsaIndex := int(posInSign / amsaSize)
	mapped := (baseSignIndex + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 12
	return VargaResult{"D12 Dvadasamsa", 12, mapped, signs[mapped], toDMS(degInAmsa)}
}


// D16: Shodashamsa / Kalamsha — movable start from Aries, Fixed start from Leo, Dual start from Saggitarius 
func D16(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 16
	amsaIndex := int(posInSign / amsaSize)

	var mapped int

	switch baseSignIndex { //0 based
	case 0,3,6,9: //movable
		mapped = 0 + amsaIndex //start from Aries
	case 1,4,7,10: //fixed
		mapped = 4 + amsaIndex //start from Leo
	case 2,5,8,11: //dual
		mapped = 8 + amsaIndex
	}

	mapped = mapped  % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 16
	fmt.Printf(" D16 Shodamsa/Kalamsha  baseSignIndex: %d , posInSign: %f amsaSize: %f mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, mapped, degInAmsa)
	return VargaResult{"D16 Shodashamsa", 16, mapped, signs[mapped], toDMS(degInAmsa)}
}

// D20: Vimshamsa 
//TODO - check why movable and dual start are exchanged
func D20(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 20
	amsaIndex := int(posInSign / amsaSize)

	var mapped int

	switch baseSignIndex { //0 based
	case 0,3,6,9: //movable
		mapped = 0 + amsaIndex //start from Aries
	case 1,4,7,10: //fixed
		mapped = 8 + amsaIndex //start from Saggi -weird
	case 2,5,8,11: //start from Leo
		mapped = 4 + amsaIndex
	}

	mapped = mapped % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 20
	fmt.Printf("D20 Vimshamsha  baseSignIndex: %d , posInSign: %f amsaSize: %f mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, mapped, degInAmsa)
	return VargaResult{"D20 Vimshamsa", 20, mapped, signs[mapped], toDMS(degInAmsa)}
}



// D24: Chaturvimshamsha 
func D24(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 24
	amsaIndex := int(posInSign / amsaSize)

	var start int
	var mapped int
	if baseSignIndex%2 != 0 { // even signs because baseSignIndex is 0 based
		start = 3
	} else {
		start = 4
	}

	mapped = start + amsaIndex
	mapped = mapped % 12 
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 24
	fmt.Printf("D24 Chaturvimshamsha  baseSignIndex: %d , posInSign: %f amsaSize: %f amsaIndex: %d mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, amsaIndex, mapped, degInAmsa)
	return VargaResult{"D24 Chaturvimshamsa", 24, mapped, signs[mapped], toDMS(degInAmsa)}
}


// D27: Nakshatramsa / Saptavimsamsa / Bhamsa 
//TODO check exchange of Cancer and Capricorn start
//FIXME Ge - 11 case works, Sc 19 fails on % change 11/12
func D27HD(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 27
	amsaIndex := int(posInSign / amsaSize)

	var start int
	var mapped int

	switch baseSignIndex {
		case 0,4,8: //fiery
			start = 0
		case 1,5,9: //earthy
			start = 3 
		case 2,6,10: //airy
			start = 6
		case 3,7,11: //watery
			start = 9
	}

	mapped = start + amsaIndex
	mapped = mapped % 12 
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 27
	fmt.Printf("D27 HD Nakshatramsa  baseSignIndex: %d , posInSign: %f amsaSize: %f amsaIndex: %d mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, amsaIndex, mapped, degInAmsa)
	return VargaResult{"D27 HD Nakshatramsa", 27, mapped, signs[mapped], toDMS(degInAmsa)}
}


// D30: Trimsamsa 
//TODO to be tested
//TODO check why there are uneven divisions of amsas of 5,6,7 degrees
func D30HD(baseSignIndex int, posInSign float64) VargaResult {
	var amsaSize int 
	var amsaIndex int

	//var start int
	var mapped int
	var doneDegrees float64 
	var posInAmsa float64
	var degInAmsa float64

	if baseSignIndex%2 == 0 { //odd sign because 0 based
		if posInSign <= 5 {
			mapped = 0 //Aries
			amsaSize = 5
			doneDegrees = 0
		} else if posInSign > 5 && posInSign <= 10 {
			mapped = 10 //Aquarius
			amsaSize = 5
			doneDegrees = 5
		} else if posInSign > 10 && posInSign <= 18 {
			mapped = 8 // Sagittarius
			amsaSize = 8
			doneDegrees = 10
		} else if posInSign > 18 && posInSign <= 25 {
			mapped = 2 // Gemini
			amsaSize = 7
			doneDegrees = 18
		} else if posInSign > 25 && posInSign <= 30 {
			mapped = 6 // Libra
			amsaSize = 5
			doneDegrees = 25
		}	
	} else { // even
		if posInSign <= 5 {
			mapped = 1 // Taurus
			amsaSize = 5
			doneDegrees = 0
		} else if posInSign > 5 && posInSign <= 12 {
			mapped = 5 // Virgo
			amsaSize = 7
			doneDegrees = 5
		} else if posInSign > 12 && posInSign <= 20 {
			mapped = 11 // Pisces
			amsaSize = 8
			doneDegrees = 12
		} else if posInSign > 20 && posInSign <= 25 {
			mapped = 9 // Capricorn
			amsaSize = 5
			doneDegrees = 20
		} else if posInSign > 25 && posInSign <= 30 {
			mapped = 7 // Scorpio
			amsaSize = 5
			doneDegrees = 25
		}
	}

	posInAmsa = posInSign - doneDegrees
	degInAmsa = posInAmsa * float64(amsaSize)
	//amsaIndex := int(posInSign / amsaSize)

//	degInAmsa := ( (30 - doneDegrees)      posInSign - float64(amsaIndex)*amsaSize) * amsaSize 
//	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * amsaSize 
	fmt.Printf("D30 HD Trimsamsa  baseSignIndex: %d , posInSign: %f amsaSize: %f amsaIndex: %d mapped: %d degInAmsa %f \n",baseSignIndex, posInSign, amsaSize, amsaIndex, mapped, degInAmsa)
	return VargaResult{"D30 HD Trimsamsa", 30, mapped, signs[mapped], toDMS(degInAmsa)}
}


func getRashiName(decimalDegree float64) string {

	fmt.Printf("decimalDegree: %f \n", decimalDegree)
	degree := math.Mod(decimalDegree, 360)
	if degree < 0 {
		degree += 360
	}
	

	fmt.Printf("degree: %f \n", degree)
	// 2. Iterate and check boundaries
	for _, p := range Rashis {
		// Standard range: Start is less than End
		if p.Start < p.End {
			if degree >= p.Start && degree < p.End {
				return p.Name
			}
		} else {
			// Wrap-around range: e.g., Start 350 to End 10
			if degree >= p.Start || degree < p.End {
				return p.Name
			}
		}
	}
	return "Unknown Rashi"
}


// ---------------- Main ----------------
//called from sh file with degrees from 0 to 359

func main() {
	// CLI flags
	deg := flag.Int("deg", 0, "Degrees (0–359)")
	min := flag.Int("min", 0, "Minutes (0–59)")
	sec := flag.Int("sec", 0, "Seconds (0–59)")
	flag.Parse()

	fmt.Printf("deg: %d  min %d  sec %d \n", *deg, *min, *sec)

	// Input validation
	if *deg < 0 || *deg >= 360 {
		fmt.Fprintln(os.Stderr, "Error: degrees must be in [0, 359]")
		os.Exit(1)
	}

	if *min < 0 || *min >= 60 || *sec < 0 || *sec >= 60 {
		fmt.Fprintln(os.Stderr, "Error: minutes/seconds must be in [0, 59]")
		os.Exit(1)
	}

	// Convert to longitude
	longitude := float64(*deg) + float64(*min)/60.0 + float64(*sec)/3600.0
	fmt.Printf("Longitude: %.4f° \n", longitude)

	//because of int casting returns the floor and hence the house number
	//baseSignIndex := (int(longitude) / 30) + 1
	baseSignIndex := int(longitude) / 30

	rashiName := getRashiName(longitude)
	fmt.Printf("Rashi : %s\n", rashiName )

	//subtracting from the house number multiplied by 30 gives the longitude within the house
	posInSign := longitude - float64(baseSignIndex*30)

	fmt.Printf("Longitude: %.4f° => Base sign: %s, Position in sign: %.4f°\n\n",
		longitude, signs[baseSignIndex], posInSign)

	// Array of calls
	funcs := []func(int, float64) VargaResult{
		D1,
		D3,
		D4HD,
		D5HD,
		D6HD,
		D7,
		D7HD,
		D8HD,
		D9,
		D9HD,
		// D9_Parasara,
		D10,
		D11HD,
		D12,
		D16,
		D20,
		D24,
		D27HD,
		D30HD,
	}

/*
	names := []string{"D1", "D3", "D7", "D9", "D10", "D12", "D16", "D20", "D24", "D27", "D30", "D40", "D45", "D60"}
*/
	fmt.Println("Varga Calculations:")
	//for i, f := range funcs {
//	var results []VargaResult
	for _, f := range funcs {
		r := f(baseSignIndex, posInSign)
		vargas = append(vargas, r)
		//fmt.Println(i, ":",  r.Varga, "  => ", r.Sign, r.DMS.String())
		fmt.Println( r.Varga, "  => ", r.Sign, r.DMS.String())
	}

	fmt.Println('\n',vargas)
}

