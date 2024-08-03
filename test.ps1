# Ensure you have the environment variables set
$USAJOBS_EMAIL = $env:USAJOBS_EMAIL
$USAJOBS_API_KEY = $env:USAJOBS_API_KEY

# Use Invoke-WebRequest to perform the same HTTP request
$response = Invoke-WebRequest -Uri "https://data.usajobs.gov/api/search" -Headers @{
        "Host" = "data.usajobs.gov"
        "User-Agent" = $USAJOBS_EMAIL
        "Authorization-Key" = $USAJOBS_API_KEY
    }

# Display the content of the response
$response.Content