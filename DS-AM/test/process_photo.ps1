# PowerShell-скрипт для тестирования обработки фотографий в фотосервисе

# Проверка наличия обязательного параметра (ID фотографии)
param (
    [Parameter(Mandatory=$true)][string]$PhotoId,
    [Parameter(Mandatory=$false)][string]$OutputFile = "processed_image.jpg"
)

# Параметры обработки
$Width = 400
$Height = 300
$Quality = 85

# URL сервиса
$ServiceUrl = "http://localhost:8081"

# Проверка статуса сервиса перед выполнением запроса
Write-Host "Проверка статуса сервиса..."
try {
    $statusResponse = Invoke-WebRequest -Uri "$ServiceUrl/status" -Method GET -ErrorAction Stop
    $statusCode = $statusResponse.StatusCode
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
}

if ($statusCode -ne 200) {
    Write-Host "Ошибка: Сервис недоступен (статус код: $statusCode)"
    Write-Host "Проверьте, запущен ли сервис и доступен ли он по адресу $ServiceUrl"
    exit 1
}

Write-Host "Сервис доступен. Выполняется запрос на обработку фотографии с ID: $PhotoId"
Write-Host "Параметры обработки: ширина=$Width, высота=$Height, качество=$Quality"
Write-Host "Файл будет сохранен как: $OutputFile"

# Создаем тело запроса в формате JSON
$RequestBody = @{
    resize = $true
    width = $Width
    height = $Height
    quality = $Quality
} | ConvertTo-Json

# Выполнение запроса к API для обработки фотографии
try {
    Invoke-WebRequest -Uri "$ServiceUrl/photos/process?id=$PhotoId" -Method POST -ContentType "application/json" -Body $RequestBody -OutFile $OutputFile -ErrorAction Stop
    
    # Проверка результата
    if (Test-Path $OutputFile) {
        $fileInfo = Get-Item $OutputFile
        
        if ($fileInfo.Length -gt 0) {
            Write-Host "Успешно! Обработанное изображение сохранено в файл: $OutputFile"
            Write-Host "Размер файла: $([math]::Round($fileInfo.Length / 1KB, 2)) KB"
            
            # Проверка содержимого файла
            $fileContent = Get-Content $OutputFile -Raw -Encoding Byte
            $isBinary = $false
            
            # Проверяем, является ли файл изображением или текстом
            foreach ($byte in $fileContent[0..100]) {
                if ($byte -lt 32 -and $byte -ne 9 -and $byte -ne 10 -and $byte -ne 13) {
                    $isBinary = $true
                    break
                }
            }
            
            if ($isBinary) {
                Write-Host "Файл содержит бинарные данные, вероятно это изображение."
            } else {
                Write-Host "Предупреждение: Полученный файл может не быть изображением. Содержимое может быть сообщением об ошибке:"
                Get-Content $OutputFile -Encoding UTF8
            }
        } else {
            Write-Host "Ошибка: Полученный файл пустой"
        }
    } else {
        Write-Host "Ошибка: Не удалось создать файл $OutputFile"
    }
} catch {
    Write-Host "Ошибка при выполнении запроса: $_"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Ответ сервера: $responseBody"
        $reader.Close()
    }
    exit 1
} 