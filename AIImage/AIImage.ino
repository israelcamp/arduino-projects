#include "tft_setup.h"
#include "mywifi.h"

#define MAX_IMAGE_WIDTH 240 // Adjust for your images

int16_t xpos = 0;
int16_t ypos = 0;

#include "SPI.h"
#include <TFT_eSPI.h>

TFT_eSPI tft = TFT_eSPI();

void setup(){
  Serial.begin(115200);

  tft.begin();
  tft.fillScreen(TFT_BLACK);

  Serial.println("\r\nInitialisation done.");
  connectToStrongestWiFi();
}

void loop() {
  delay(10000);
}