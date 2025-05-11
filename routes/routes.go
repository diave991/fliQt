package routes

import (
	"fliQt/controllers"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, empCtrl *controllers.EmployeeController, leaveCtrl *controllers.LeaveController, attCtrl *controllers.AttendanceController) {
	api := r.Group("/api/v1")
	{
		// Employee
		emp := api.Group("/employees")
		{
			emp.POST("", empCtrl.Create)
			emp.GET("", empCtrl.GetAll)
			emp.GET("/:id", empCtrl.GetByID)
			emp.PUT("/:id", empCtrl.Update)
			emp.DELETE("/:id", empCtrl.Delete)
		}
		// Leave
		lv := api.Group("/leaves")
		{
			lv.POST("", leaveCtrl.Create)
			lv.GET("", leaveCtrl.GetAll)
		}
		// Attendance
		at := api.Group("/attendance")
		{
			at.POST("", attCtrl.Create)
			at.GET("", attCtrl.GetAll)
			at.GET("/:id", attCtrl.GetByID)
			at.GET("/by_employee", attCtrl.GetByEmployee) // ?employee_id=xx[&date=YYYY-MM-DD]
		}
	}
}
