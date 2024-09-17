# GIS Feature Collection API

## Overview

This API application provides a feature collection for a GIS rendering application. It enables users to manage geolocations and retrieve them in GeoJSON format. The API supports the following functionalities:

- Create Geolocations: Add places with their descriptions and coordinates (latitude/longitude).
- Retrieve Geolocations: Fetch a feature collection of stored geolocations in GeoJSON format.
- Caching: Implemented caching mechanism to enhance performance.
- Data Storage: Utilizes PostgreSQL for data storage.

## Application Architecture

The application architecture is inspired by the [Hexagonal Architecture pattern](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749), featuring a structured separation of concerns across different layers:


### 1. Domain Layer

- **Entities**: Represent the core business objects or models of the application. They encapsulate the business logic and rules. For instance, in a GIS application, an Entity could be a `Geolocation` with methods for manipulating its data. This layer is purely focused on the business logic and does not concern itself with how data is stored or provided.
  
- **Ports**: Define the interfaces for interacting with the core business logic. These interfaces specify the operations that can be performed on the domain entities without specifying the details of data storage or retrieval. For example, a `GeolocationService` interface might define methods such as `CreateGeolocations` and `GetGeolocationsGeoJSON`.

### 2. Application Layer

- **Services**: Implement the business logic and use the domain entities and ports. They interact with the ports to execute the application's use cases. For example, a `GeolocationService` might implement methods for creating a new geolocation, validating it, and saving it using the port interfaces. This layer coordinates the application's activities and ensures the proper execution of business rules.

### 3. Adapter Layer

- **Repositories**: Implement the data access logic and interact with external data sources like databases or external APIs. They adapt the domain's ports to specific technologies, handling the persistence and retrieval of data. For example, a `GeolocationRepository` could manage saving and retrieving geolocations from a PostgreSQL database.

- **Handlers**: Implement the interfaces that connect the application to external inputs or outputs, such as web requests or messaging systems. They adapt the inputs to the application's needs, translating HTTP requests or other inputs into calls to the application services. For instance, a `GeolocationHandler` might process HTTP requests and use the `GeolocationService` to perform operations like creating or retrieving geolocations.


## Technologies Used

- Go: Programming language used for implementing the API.
- PostgreSQL: Database for storing geolocation data.
- Redis: Caching mechanism to improve performance.

## Running the Application

### Setup

1. Install and Configure PostgreSQL and Redis:
   - Follow the installation guides for [PostgreSQL](https://www.postgresql.org/download/) and [Redis](https://redis.io/download/).
   - Ensure both services are running before starting the application.

2. Create Configuration File:
   - Inside the `conf` folder, create a configuration file based on `example.yaml`. For example, create `development.yaml` with the appropriate settings.

### Migration

1. Install [Go Migrate](https://github.com/golang-migrate/migrate)
2. Set Environment Variables:
   - Export your PostgreSQL URL:

     ```bash 
     export POSTGRESQL_URL='postgres://{name}:{password}@{host}:{port}/{db-name}?sslmode=disable' 
     ```
3. Run Migrations:
   - Execute the migration command:
     ```bash 
     migrate -database ${POSTGRESQL_URL} -path ./migration up` 
     ```

4. Manual Migration (if needed):
   - Copy the SQL files from the `migration` directory and execute them directly in your database.


### Build and Run

1. **Build the Application**:
   - Build the application binary:
     ```bash
     go build -o cmd/main cmd/main.go
     ```

2. **Run the Application**:
   - Execute the built application:
     ```bash
     cmd/main
     ```

3. **Run in Development Mode**:
   - Run the application in development mode:
     ```bash
     go run cmd/main.go
     ```
   - Alternatively, you can use [nodemon](https://www.npmjs.com/package/nodemon) for automatic reloading:
     ```bash
     nodemon --exec go run ./cmd/main.go --signal SIGTERM
     ```
## API Endpoints

### Create Geolocations

- Endpoint: `POST /api/create-batch`
- Request Body: 
    ```json 
        { 
            "items": [ 
                { 
                    "name": "Bandar Udara Internasional Halim Perdanakusuma", 
                    "description": "Bandara Internasional", 
                    "lat": -6.265135263342884, 
                    "lng": 106.88583174452127 
                } 
            ] 
        }
- Response:
    ```json
        { "message": "Geolocation created successfully" }
     ```

### Retrieve Geolocations

- Endpoint: `GET /api/map-data`
- Response:
    ```json
        {
            "type": "FeatureCollection",
            "features": [
                {
                    "type": "Feature",
                    "geometry": {
                        "type": "Point",
                        "coordinates": [
                            106.88583174452127,
                            -6.265135263342884
                        ]
                    },
                    "properties": {
                        "description": "Bandara Internasional",
                        "name": "Bandar Udara Internasional Halim Perdanakusuma"
                    }
                }
            ]
        }
     ```
- Explanation:

    You can download this response as a JSON file using tools like Postman.
    Upload the downloaded JSON file to [Kepler.gl](https://kepler.gl/demo ) to visualize the geolocations on a map. Kepler.gl provides a user-friendly interface for visualizing geospatial data and can help you analyze and present the geolocations effectively.
## License

[MIT License](LICENSE) 

Feel free to adjust any part of this draft to better fit your application or preferences!