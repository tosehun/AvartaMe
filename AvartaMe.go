//
// Avarta Me
// This is one example of "github.com/GoBootcamp/avatarme"
//
// By tosehun@gmail.com
//
// MYCINDY Lisence : if you know Cindy, it's free for any use.
//
// +--------+-----------+
// | color  |     map   |
// | 1 byte |   2bytes  |
// +--------+-----------+

package main

import (
    "fmt"
    "os"
    "crypto/md5"
    "image"
    "image/color"
    "image/png"
    "image/draw"
)

const SIZE = 6

func main(){

    if len(os.Args) != 2 {
        fmt.Println("Avarta Me needs one argument : " , len(os.Args)-1 , "argument(s) we have.")
        return 
    }

    inputString := os.Args[1]
    outputFileName := inputString + ".png"

    idData := getIdDataFromString(inputString)
    fmt.Println("ID : ", idData, len(idData))
    foreColor := getColorFromIdData(idData)
    fmt.Println("Color : ", foreColor)
    mapData := getMapfromIdData(idData)
    //fmt.Println("map : ", len(mapData), " : " , mapData)
    drawOnConsole(mapData);
    makePNG(mapData, foreColor, 300, outputFileName)
    fmt.Println("Created ", outputFileName)

}

//we need hashed data only 3 bytes.
// 20bit-size data allows to make small map * 64bit color enough.
func getIdDataFromString(str string) []byte {
    data := md5.Sum([]byte(str))
    return data[:3]

}

//we use just one color. just pick a color from first byte
func getColorFromIdData(idData []byte) color.Color {
    levelTable := []byte{0x33, 0x66, 0x99, 0xCC}
    r := idData[0] >> 6           //first 2 bits
    g := (idData[0] >> 4 ) & 0x3  //second 2 bits
    b := (idData[0] >> 2) & 0x3   //third 2 bits
    //not use the last 2 bits now. reserved.

	c := color.RGBA{levelTable[r],levelTable[g],levelTable[b],0xf0}
	return c
}

//create map data using 2nd~3rd byte 
func getMapfromIdData(idData []byte) [][]bool {
	mapData := make([][]bool, SIZE*SIZE)
	for i:=0;i<SIZE;i++ {
	    mapData[i] = make([]bool, SIZE)
	}

	//available 14 positions (left side)
    mapData[1][0] = (idData[1] & 0x80) > 0 //1st byte - 1st bit
    mapData[2][0] = (idData[1] & 0x40) > 0 //2nd bit
    mapData[1][1] = (idData[1] & 0x20) > 0 //3rd bit
    mapData[2][1] = (idData[1] & 0x10) > 0 //4th bit
    mapData[1][2] = (idData[1] & 0x08) > 0 //5th bit
    mapData[2][2] = (idData[1] & 0x04) > 0 //6th bit
    mapData[0][3] = (idData[1] & 0x02) > 0 //7th bit
    mapData[1][3] = (idData[1] & 0x01) > 0 //8th bit
    mapData[0][4] = (idData[2] & 0x80) > 0 //2nd byte - 1th bit
    mapData[1][4] = (idData[2] & 0x40) > 0 //2nd bit
    mapData[2][4] = (idData[2] & 0x20) > 0 //3rd bit
    mapData[0][5] = (idData[2] & 0x10) > 0 //4th bit
    mapData[1][5] = (idData[2] & 0x08) > 0 //5th bit
    mapData[2][5] = (idData[2] & 0x04) > 0 //6th bit

    //make symmetric (to right side)
    for i:=0;i<SIZE/2;i++ {
        for j:=0;j<SIZE;j++ {
            mapData[SIZE-i-1][j] = mapData[i][j]
        }
    }

    //TODO : rotate - 180 degree
    if (idData[2] & 0x02) > 0 { //7th bit

    }

	return mapData
}

func drawOnConsole(mapData [][]bool){
    for i:=0;i<SIZE;i++ {
        for j:=0;j<SIZE;j++ {
            if mapData[j][i] == true {
        	    fmt.Printf(" *")
            }else{
            	fmt.Printf(" _")
            }
        }
        fmt.Printf("\n")
    }
}

func makePNG(mapData [][]bool, foreColor color.Color, width int, filePath string) {
	cellWidth := int(width/SIZE)
    imageRect := image.Rect(0,0,width,width)
    imageData := image.NewRGBA(imageRect)
    bgColor := color.RGBA{0xee,0xee,0xee,0xff}

    draw.Draw(imageData, 
        imageData.Bounds(), 
        &image.Uniform{bgColor},
        image.ZP,
        draw.Src)

    for i:=0;i<SIZE;i++ {
        for j:=0;j<SIZE;j++ {

            if(mapData[j][i] == true){
            	x := j*cellWidth
            	y := i*cellWidth
            	draw.Draw(imageData,
            		image.Rect(x, y, x+cellWidth, y+cellWidth),
                    &image.Uniform{foreColor},
                    image.ZP,
                    draw.Src)
            }

        }
    }


    outFile, err := os.Create("./"+filePath)
    if err != nil {
             fmt.Println("FileOpen Error: ", err)
             os.Exit(1)
    }

    err = png.Encode(outFile, imageData)
    if err != nil {
             fmt.Println("Encode Error: ", err)
             os.Exit(1)
    }
}