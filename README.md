# Drobo MQTT

Run this service in order to publish messages about Drobo to MQTT. Creates a
device with several sensors:
- Drobo's state
- Drobo's capacity
- Drobo's % capacity used
- Each disk's state

### Getting running

1. Run the build command
2. Copy the binary to a willing host
3. Copy the systemctl config (example/drobo-mqtt.service) to the savant host
   under /lib/systemd/system/drobo-mqtt.service
4. Run `sudo systemctl daemon-reload`
5. Enable the service via `sudo systemctl enable drobo-mqtt`
6. Start the service via `sudo systemctl start drobo-mqtt`
