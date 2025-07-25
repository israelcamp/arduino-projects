#include <U8g2lib.h>

U8G2_SSD1309_128X64_NONAME0_F_HW_I2C u8g2(U8G2_R0, /* reset=*/ U8X8_PIN_NONE);

const int trigPin = 2;
const int echoPin = 3;
const int buzzerPin = 7;
const bool useBuzzer = true; // change this if buzzer not available
const int fontSize = 14;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
  pinMode(trigPin, OUTPUT);
  pinMode(echoPin, INPUT);
  u8g2.begin();
  u8g2.setFont(u8g2_font_ncenB14_tr);

  pinMode(buzzerPin, OUTPUT);
}

void loop() {
  // put your main code here, to run repeatedly:
  digitalWrite(trigPin, LOW);
  delayMicroseconds(2);
  digitalWrite(trigPin, HIGH);
  delayMicroseconds(10);
  digitalWrite(trigPin, LOW);

  long duration = pulseIn(echoPin, HIGH); // in microseconds

  float distanceCm = (duration * 0.0343) / 2.0;

  // --- convert float to string ---
  char buf[8];
  dtostrf(distanceCm, 4, 2, buf);  // e.g. "12.3" :contentReference[oaicite:4]{index=4}

  // --- draw title ---
  const char *title = "DISTANCE";
  // u8g2.setFont(u8g2_font_fub42_tr);            // big font :contentReference[oaicite:5]{index=5}
  int16_t w = u8g2.getStrWidth(title);         // measure width :contentReference[oaicite:6]{index=6}
  u8g2.clearBuffer();
  u8g2.drawStr((128 - w) / 2, fontSize, title);      // vertical pos “48” near top

  // --- draw value + units ---
  int16_t wVal = u8g2.getStrWidth(buf);
  int16_t wUnit = u8g2.getStrWidth(" cm");
  int16_t x0 = (128 - (wVal + wUnit)) / 2;
  int16_t y0 = 2 * fontSize + 6;
  u8g2.drawStr(x0, y0, buf);                   // e.g. "12.3"
  u8g2.drawStr(x0 + wVal, y0, " cm"); 

  if (useCBuzzer && distancem < 10) {
    tone(buzzerPin, 2000, 1000);
    const char *clearOff = "CAI FORA";
    int16_t w1 = u8g2.getStrWidth(clearOff);         // measure width :contentReference[oaicite:6]{index=6}
    u8g2.drawStr((128 - w1) / 2, 3 * fontSize + 12, clearOff);
  }

  u8g2.sendBuffer();

  delay(500);
}
