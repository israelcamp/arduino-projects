#include "secrets.h"
#include <WiFi.h>

int connectAndReturnStrength(const char *ssid, const char *password) {
  WiFi.begin(ssid, password);
  WiFi.setSleep(false);

  int tries = 0;
  Serial.printf("Trying to connect to %s", ssid);
  while (WiFi.status() != WL_CONNECTED && tries < 20) {
    delay(500);
    Serial.print(".");
    tries++;
  }
  Serial.println("");
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("Could not connect");
    return -1000;
  }

  int strength = WiFi.RSSI();
  Serial.printf("Connected to %s with strength %d \n", ssid, strength);
  return strength;
}

void connectToStrongestWiFi() {
  const int wifi1_strength = connectAndReturnStrength(WiFi1_ssid, WiFi1_password);
  // const int wifi2_strength = connectAndReturnStrength(WiFi2_ssid, WiFi2_password);
  // if (wifi1_strength > wifi2_strength) {
  //   connectAndReturnStrength(WiFi1_ssid, WiFi1_password);
  // }
}