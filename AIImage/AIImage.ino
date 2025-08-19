#include "tft_setup.h"
#include "mywifi.h"
#include <HTTPClient.h>
#include <TJpg_Decoder.h>
#include <SPI.h>
#include <TFT_eSPI.h>

#include "NotoSansBold15.h"
#include "NotoSansBold36.h"

// The font names are arrays references, thus must NOT be in quotes ""
#define AA_FONT_SMALL NotoSansBold15
#define AA_FONT_LARGE NotoSansBold36

TFT_eSPI tft = TFT_eSPI();

const int16_t fontHeight = 36;
const int noPersonLedPin = 33;
const int yesPersonLedPin = 12;
const char* headers[] = {"HasPerson"};
String hasPerson = "no";

bool tftOutput(int16_t x,int16_t y,uint16_t w,uint16_t h,uint16_t *bmp){
  if (y >= tft.height()) return 0;
  tft.pushImage(x, y, w, h, bmp);
  return 1;
}

void setup(){
  Serial.begin(115200);

  tft.begin();
  tft.setTextColor(0xFFFF, 0x0000);
  tft.fillScreen(TFT_BLACK);

  TJpgDec.setJpgScale(1);
  TJpgDec.setSwapBytes(true);
  TJpgDec.setCallback(tftOutput);   

  tft.setTextColor(TFT_WHITE, TFT_BLACK);
  tft.loadFont(AA_FONT_LARGE);

  pinMode(TFT_BL, OUTPUT);
  pinMode(noPersonLedPin, OUTPUT);
  pinMode(yesPersonLedPin, OUTPUT);

  digitalWrite(TFT_BL, HIGH);
  
  // HELLO
  const char* helloText = "HELLO!";
  tft.setTextDatum(MC_DATUM);
  tft.drawString(helloText, tft.width()/2, tft.height()/2);
  delay(2000);

  tft.fillScreen(TFT_BLACK);
  tft.setTextDatum(MC_DATUM);
  tft.drawString("System", tft.width()/2, tft.height()/2 - 30);
  tft.drawString("by Israel", tft.width()/2, tft.height()/2);
  delay(2000);

  tft.fillScreen(TFT_BLACK);
  const char* wifiText = "Connecting...";
  tft.setTextDatum(MC_DATUM);
  tft.drawString(wifiText, tft.width()/2, tft.height()/2);
  delay(2000);

  connectToStrongestWiFi();

  tft.fillScreen(TFT_BLACK);
  const char* ready = "ENJOY!!";
  tft.setTextDatum(MC_DATUM);
  tft.drawString(ready, tft.width()/2, tft.height()/2);
  delay(2000);

  tft.fillScreen(TFT_BLACK);
}

void loop(){
  HTTPClient http;
  http.begin(serverUrl);
  http.collectHeaders(headers, 1);

  if (http.GET()==HTTP_CODE_OK){
    int len = http.getSize();
    auto *buf = (uint8_t*) heap_caps_malloc(len, MALLOC_CAP_INTERNAL|MALLOC_CAP_8BIT);
    http.getStreamPtr()->readBytes(buf,len);

    uint16_t w,h; TJpgDec.getJpgSize(&w,&h,buf,len);

    TJpgDec.drawJpg(0, 0, buf, len);

    free(buf);

    hasPerson = http.header("HasPerson");

  } else {
    tft.fillScreen(TFT_BLACK);
  }
  http.end();

  if (hasPerson == "yes") {
    digitalWrite(yesPersonLedPin, HIGH);
    digitalWrite(noPersonLedPin, LOW);
  } else {
    digitalWrite(yesPersonLedPin, LOW);
    digitalWrite(noPersonLedPin, HIGH);
  }

  // Wait before drawing again
  delay(250);

}

