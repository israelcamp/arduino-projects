#include "tft_setup.h"
#include "mywifi.h"
#include <HTTPClient.h>
#include <TJpg_Decoder.h>
#include "SPI.h"
#include <TFT_eSPI.h>

TFT_eSPI tft = TFT_eSPI();

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

  // The jpeg image can be scaled by a factor of 1, 2, 4, or 8
  TJpgDec.setJpgScale(1);

  // The byte order can be swapped (set true for TFT_eSPI)
  TJpgDec.setSwapBytes(true);

  // The decoder must be given the exact name of the rendering function above
  TJpgDec.setCallback(tftOutput);   


  connectToStrongestWiFi();

}

void loop(){
  HTTPClient http;
  http.begin("http://192.168.15.12:8090/aicapture");
  if (http.GET()==HTTP_CODE_OK){
    int len = http.getSize();
    auto *buf = (uint8_t*) heap_caps_malloc(len, MALLOC_CAP_INTERNAL|MALLOC_CAP_8BIT);
    http.getStreamPtr()->readBytes(buf,len);

    uint16_t w,h; TJpgDec.getJpgSize(&w,&h,buf,len);

    TJpgDec.drawJpg(0, 0, buf, len);

    free(buf);
  }
  http.end();

  // Wait before drawing again
  delay(2000);
}
