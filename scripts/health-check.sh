set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

ENVIRONMENT=${1:-development}
DOMAIN=${2:-localhost}

echo -e "${GREEN}🔍 Starting health check for environment: ${ENVIRONMENT}${NC}"

check_endpoint() {
    local url=$1
    local name=$2
    local expected_status=${3:-200}
    
    echo -n "Checking $name... "
    
    if response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null); then
        if [ "$response" = "$expected_status" ]; then
            echo -e "${GREEN}✅ OK (HTTP $response)${NC}"
            return 0
        else
            echo -e "${RED}❌ FAILED (HTTP $response, expected $expected_status)${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ FAILED (Connection error)${NC}"
        return 1
    fi
}

check_container() {
    local container_name=$1
    local service_name=$2
    
    echo -n "Checking $service_name container... "
    
    if docker ps --filter "name=${container_name}" --filter "health=healthy" | grep -q "${container_name}"; then
        echo -e "${GREEN}✅ Healthy${NC}"
        return 0
    else
        echo -e "${RED}❌ Unhealthy${NC}"
        return 1
    fi
}

echo -e "${YELLOW}🐳 Checking Docker containers...${NC}"
services=("postgres" "redis" "backend" "frontend" "nginx")
failed_containers=0

for service in "${services[@]}"; do
    if [ "$ENVIRONMENT" = "production" ]; then
        container_name="ecommerce-${service}-prod"
    else
        container_name="ecommerce-${service}"
    fi
    
    if ! check_container "$container_name" "$service"; then
        ((failed_containers++))
    fi
done

echo -e "${YELLOW}🌐 Checking HTTP endpoints...${NC}"
failed_endpoints=0

if [ "$ENVIRONMENT" = "production" ]; then
    protocol="https"
    base_url="https://${DOMAIN}"
else
    protocol="http"
    base_url="http://localhost"
fi

if ! check_endpoint "${base_url}/health" "Health endpoint"; then
    ((failed_endpoints++))
fi

if ! check_endpoint "${base_url}/" "Frontend"; then
    ((failed_endpoints++))
fi

if ! check_endpoint "${base_url}/api/health" "Backend API"; then
    ((failed_endpoints++))
fi

if [ "$ENVIRONMENT" = "production" ]; then
    if ! check_endpoint "${base_url}:3001" "Grafana"; then
        ((failed_endpoints++))
    fi
    
    if ! check_endpoint "${base_url}:9090" "Prometheus"; then
        ((failed_endpoints++))
    fi
fi

echo -e "${YELLOW}📊 Health Check Summary:${NC}"
echo -e "  Containers: $((5 - failed_containers))/5 healthy"
echo -e "  Endpoints: $((3 - failed_endpoints))/3 responding"

if [ $failed_containers -eq 0 ] && [ $failed_endpoints -eq 0 ]; then
    echo -e "${GREEN}🎉 All systems are healthy!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some systems are not healthy. Please check the logs.${NC}"
    exit 1
fi
