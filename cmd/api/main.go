package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"template-fiber-v3/configs"
	server "template-fiber-v3/internal/infrastructure"
	"time"
)

// ฟังก์ชันสำหรับปิดเซิร์ฟเวอร์อย่างปลอดภัย (Graceful Shutdown)
func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// สร้าง Context ที่ใช้สำหรับรอฟังสัญญาณจากระบบปฏิบัติการ
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // เรียกใช้ stop() เพื่อยกเลิก Context เมื่อเสร็จสิ้น

	// รอรับสัญญาณจากระบบปฏิบัติการ
	<-ctx.Done()

	log.Println("กำลังปิดเซิร์ฟเวอร์อย่างปลอดภัย, กด Ctrl+C อีกครั้งหากต้องการปิดทันที")

	// กำหนด timeout 10 วินาทีให้เซิร์ฟเวอร์สามารถจัดการคำขอที่กำลังดำเนินการอยู่ก่อนปิด
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// เรียกใช้ฟังก์ชัน ShutdownWithContext เพื่อปิดเซิร์ฟเวอร์อย่างปลอดภัย
	if err := fiberServer.ShutdownWithContext(shutdownCtx); err != nil {
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Println("ไม่สามารถปิดเซิร์ฟเวอร์ได้ภายในเวลาที่กำหนด กำลังปิดทันที")
		} else {
			log.Printf("ไม่สามารถปิดเซิร์ฟเวอร์ได้: %v", err)
		}
	}

	log.Println("เซิร์ฟเวอร์ปิดแล้ว")

	// แจ้งให้ Goroutine หลักทราบว่ากระบวนการปิดเซิร์ฟเวอร์เสร็จสิ้น
	done <- true
}

func main() {

	// โหลด config
	configs, err := configs.LoadConfig()
	if err != nil {
		log.Fatalln("Load ENV ERROR : ", err.Error())
	}

	// สร้างเซิร์ฟเวอร์
	server := server.New(configs)

	// ลงทะเบียนเส้นทาง API ต่าง ๆ
	server.RegisterFiberRoutes()

	// สร้าง channel สำหรับตรวจสอบสถานะการปิดเซิร์ฟเวอร์
	done := make(chan bool, 1)

	// เริ่มต้นเซิร์ฟเวอร์ใน Goroutine แยก
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("พบ panic ในเซิร์ฟเวอร์: %v", r)
				done <- true
			}
		}()

		err = server.Listen(fmt.Sprintf(":%d", configs.Port))
		if err != nil {
			panic(fmt.Sprintf("เกิดข้อผิดพลาดในเซิร์ฟเวอร์: %s", err))
		}
	}()

	// เรียกใช้ฟังก์ชันปิดเซิร์ฟเวอร์ใน Goroutine แยก
	go gracefulShutdown(server, done)

	// รอจนกว่ากระบวนการปิดเซิร์ฟเวอร์จะเสร็จสมบูรณ์
	<-done
	log.Println("กระบวนการปิดเซิร์ฟเวอร์เสร็จสมบูรณ์")
}
