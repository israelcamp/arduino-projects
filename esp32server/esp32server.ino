#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include <WiFi.h>
#include <WebServer.h>
#include "secrets.h"

WebServer server(82);

const char *idleText = "Nada acontece...";
const char *badText = "PESSOA ENCONTRADA";
bool isIdle = true;

Adafruit_SSD1306 display(128, 64, &Wire, -1);

void handleRoot() {
  isIdle = false;
  Serial.println("Received request");
  server.send(200, "text/plain", "Hi from ESP32 sever!");
}

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  // SDA = IO15, SCL = IO13, 100 kHz
  Wire.begin(15, 13, 100000);

  // connecting to wifi
  connectToStrongestWiFi();

  Serial.print("Camera Ready! Use 'http://");
  Serial.print(WiFi.localIP());
  Serial.println("' to connect");

  server.on("/", handleRoot);
  server.begin();

  if (!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) {
    Serial.println("SSD1306 allocation failed");
    for (;;);
  }

  
}

void loop() {
  // put your main code here, to run repeatedly:
  server.handleClient();

  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SSD1306_WHITE);
  display.setCursor(0, 0);
  display.println(isIdle ? idleText : badText);
  display.display();

  delay(100);
}

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
  IPAddress local_IP(192, 168, 15, 119);
  IPAddress gateway(192, 168, 15, 1);
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
  IPAddress local_IP(192, 168, 0, 109);
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