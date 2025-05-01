package integration

import (
	"bytes"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestLoadBalancerWithAB(t *testing.T) {
	port := "8080"

	resp, err := http.Get("http://localhost:" + port + "/health")
	if err == nil {
		resp.Body.Close()
		t.Fatalf("Порт %s уже занят другим процессом!", port)
	}

	cmd := exec.Command("go", "run", "../../cmd/main.go", "-port", port)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		t.Fatalf("Ошибка запуска балансировщика: %v", err)
	}
	defer cmd.Process.Kill()

	t.Log("Ожидание запуска балансировщика...")
	ready := false
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://localhost:" + port + "/health")
		if err == nil {
			resp.Body.Close()
			ready = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	if !ready {
		t.Fatalf("Балансировщик не запустился. Вывод:\n%s", out.String())
	}

	t.Log("Запуск теста с умеренной нагрузкой (500 запросов, 100 одновременных)")
	abCmd := exec.Command("ab", "-n", "5000", "-c", "100", "http://localhost:"+port+"/health")
	if output, err := abCmd.CombinedOutput(); err != nil {
		t.Logf("Тест с умеренной нагрузкой завершился с ошибкой: %v\nВывод:\n%s", err, output)
	}

	t.Log("Запуск теста с высокой нагрузкой (1000 запросов, 1000 одновременных)")
	abCmd = exec.Command("ab", "-n", "1000", "-c", "1000", "http://localhost:"+port+"/health")
	if output, err := abCmd.CombinedOutput(); err != nil {
		t.Errorf("Тест с высокой нагрузкой завершился с ошибкой: %v\nВывод:\n%s", err, output)
	}

	t.Logf("Вывод балансировщика:\n%s", out.String())
}
