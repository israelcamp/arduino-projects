const int sensorPin = 8;

int count = 0;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);

  pinMode(sensorPin, INPUT);

}

void loop() {
  // put your main code here, to run repeatedly:
  const bool isHigh = digitalRead(sensorPin);
  if (isHigh) {
    count++;
    Serial.print(count * 50);
    Serial.println(" SOMETHING!");
  } else {
    count = 0;
  }
  delay(50);
}
