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
	return VargaResult{"D1 Rasi", baseSignIndex, signs[baseSignIndex], toDMS(posInSign)}
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
	return VargaResult{"D3 Drekkana", mapped, signs[mapped], toDMS(degInAmsa)}
}

//D4 HD as per PVRN
func D4HD(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 7.5 
	amsaIndex := int(posInSign / amsaSize)
	offsets := []int{0, 3, 6, 9}
	mapped := (baseSignIndex + offsets[amsaIndex]) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 3
//	fmt.Println("baseSignIndex: ", baseSignIndex, " posInSign: ", posInSign,  " amsaIndex:", amsaIndex, " mapped: ", mapped, " degInAmsa: ", degInAmsa)
	return VargaResult{"D4 HD Chaturthamsa", mapped, signs[mapped], toDMS(degInAmsa)}
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
	return VargaResult{"D5 HD Panchamsa", mapped, signs[mapped], toDMS(degInAmsa)}
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
	return VargaResult{"D6 HD Shasthamsa", mapped, signs[mapped], toDMS(degInAmsa)}
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
	return VargaResult{"D7 Saptamsa", mapped, signs[mapped], toDMS(degInAmsa)}
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
	return VargaResult{"D9 Navamsa", mapped, signs[mapped], toDMS(degInAmsa)}
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

	return VargaResult{"D9 HD", mapped, signs[mapped], toDMS(degInAmsa)}
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
func D10(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 10
	amsaIndex := int(posInSign / amsaSize)
	start := baseSignIndex
	if baseSignIndex%2 != 0 { // even
		start = (baseSignIndex + 8) % 12
	}
	mapped := (start + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 10
	return VargaResult{"D10 Dasamsa", mapped, signs[mapped], toDMS(degInAmsa)}
}

// D12: Dvadasamsa
func D12(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 12
	amsaIndex := int(posInSign / amsaSize)
	mapped := (baseSignIndex + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 12
	return VargaResult{"D12 Dvadasamsa", mapped, signs[mapped], toDMS(degInAmsa)}
}


// D16: Shodashamsa — odd signs start from Aries, even signs start from Libra
func D16(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 16
	amsaIndex := int(posInSign / amsaSize)

	start := 0 // Aries
	if baseSignIndex%2 != 0 { // even sign
		start = 6 // Libra
	}

	mapped := (start + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 16
	return VargaResult{"D16 Shodashamsa", mapped, signs[mapped], toDMS(degInAmsa)}
}

// D20: Vimshamsa — odd signs start from Leo, even signs start from Cancer
func D20(baseSignIndex int, posInSign float64) VargaResult {
	amsaSize := 30.0 / 20
	amsaIndex := int(posInSign / amsaSize)

	start := 4 // Leo
	if baseSignIndex%2 != 0 { // even sign
		start = 3 // Cancer
	}

	mapped := (start + amsaIndex) % 12
	degInAmsa := (posInSign - float64(amsaIndex)*amsaSize) * 20
	return VargaResult{"D20 Vimshamsa", mapped, signs[mapped], toDMS(degInAmsa)}
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
		D9,
		D9HD,
		// D9_Parasara,
		D10,
		D12,
		D16,
		D20,
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

