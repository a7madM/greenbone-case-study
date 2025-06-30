# Greenbone Case Study

### Prerequisites
- Docker and Docker Compose installed
- Make utility installed

### Running the Application

1. Clone the repository:
```bash
git clone https://github.com/a7madM/greenbone-case-study

cd greenbone-case-study
```

2. Start the application using Make:
```bash
make dev
```


### Makefile Usage

The project includes a Makefile with the following commands:

- `make up` - Start all services using Docker Compose
- `make down` - Stop and remove all containers
- `make logs` - View logs from all services
- `make build` - Build Docker images
- `make test` - Run test suite
- `make restart` - Stop and Restarts container

### Access

Once running, the application will be available at:
- Main application: http://localhost:3000
- API documentation: http://localhost:3000/swagger

![Screenshot 2025-06-30 at 11 57 52](https://github.com/user-attachments/assets/fd399833-95ff-4365-95ed-e73ec2029c48)

### System Overview
![Screenshot 2025-06-30 at 12 43 08](https://github.com/user-attachments/assets/73cb03f5-97a6-4410-8b0f-d366601e97c9)

### Sequence Diagram
![Screenshot 2025-06-30 at 12 40 12](https://github.com/user-attachments/assets/f3da35c4-8941-412b-9361-76be90972b80)

