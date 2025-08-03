#include <WiFi.h>
#include "secrets.h"

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

void configWiFi1() {
  // Configure static IP (adjust to your network’s settings)
  IPAddress local_IP(192, 168, 15, 117);
  IPAddress gateway(192, 168, 1, 1);
  IPAddress subnet(255, 255, 255, 0);
  IPAddress primaryDNS(8, 8, 8, 8);       // Recommended
  IPAddress secondaryDNS(8, 8, 4, 4);     // Optional

  // Apply static IP config
  if (!WiFi.config(local_IP, gateway, subnet, primaryDNS, secondaryDNS)) {
    Serial.println("STA Failed to configure");
  }
}

void configWiFi2() {
  // Configure static IP (adjust to your network’s settings)
  IPAddress local_IP(192, 168, 0, 107);
  IPAddress gateway(192, 168, 1, 1);
  IPAddress subnet(255, 255, 255, 0);
  IPAddress primaryDNS(8, 8, 8, 8);       // Recommended
  IPAddress secondaryDNS(8, 8, 4, 4);     // Optional

  // Apply static IP config
  if (!WiFi.config(local_IP, gateway, subnet, primaryDNS, secondaryDNS)) {
    Serial.println("STA Failed to configure");
  }
}

void connectToStrongestWiFi() {
  const int wifi1_strength = connectAndReturnStrength(WiFi1_ssid, WiFi1_password);
  const int wifi2_strength = connectAndReturnStrength(WiFi2_ssid, WiFi2_password);
  if (wifi1_strength > wifi2_strength) {
    configWiFi1();
    connectAndReturnStrength(WiFi1_ssid, WiFi1_password);
  } else {
    configWiFi2();
    connectAndReturnStrength(WiFi2_ssid, WiFi2_password);
  }
}