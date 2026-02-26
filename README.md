# goTimeseries

## Input

```
topic format: timeseries/#
```

The part efter / is the **bucket** in the TS database

Data format:

```json
{
   "data": { "measurement":"GoLang Project","tags": {
              "project_id":"12345"
            },
            "fields": {
              "value":42
            },
            "timestamp":"2026-01-15T20:39:24.003202+01:00"
          }
}
```
## Konfiguration fil

Filename: **config.conf**

### MQTT Konfiguration
   Egenskab        | Standardværdi      | Beskrivelse                                      |
 |-----------------|--------------------|--------------------------------------------------|
 | `mqtt_address`  | `localhost:1883`   | Adressen til MQTT-brokeren (host:port).          |
 | `mqtt_username` | `mqtt-user`        | Brugernavn til MQTT-brokeren.                    |
 | `mqtt_password` | `mqtt-password`    | Adgangskode til MQTT-brokeren.                   |

### InfluxDB Konfiguration
 | Egenskab         | Standardværdi                                                                                     | Beskrivelse                                      |
 |------------------|---------------------------------------------------------------------------------------------------|--------------------------------------------------|
 | `influxdb_url`   | `http://localhost:8086`                                                                           | URL til InfluxDB-instansen.                      |
 | `influxdb_token` | `tzhe2Ax2rtX07xyyXP_BcRtZYEftw9sCgMtS3qFnuSJ93PkFqEnRlzH1_rxst_esEwaAShMX31WDsRnz7KrTww==` | Autentificeringstoken til InfluxDB.              |
 | `influxdb_org`   | `my-org`                                                                                          | Organisationsnavn i InfluxDB.                   |

---