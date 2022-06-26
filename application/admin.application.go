package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spo-iitk/ras-backend/rc"
	"github.com/spo-iitk/ras-backend/student"
	"github.com/spo-iitk/ras-backend/util"
)

type ApplicantsByRole struct {
	StudentRCID uint   `json:"student_rc_id"`
	ResumeLink  string `json:"resume_link"`
	ProformaID  uint   `json:"proforma_id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
}

type studentAdminsideResponse struct {
	ID                           uint    `json:"id"`
	Name                         string  `json:"name"`
	Email                        string  `json:"email"`
	CPI                          float64 `json:"cpi"`
	ProgramDepartmentID          uint    `json:"program_department_id"`
	SecondaryProgramDepartmentID uint    `json:"secondary_program_department_id"`
	TenthBoard                   string  `json:"tenth_board"`
	TenthYear                    uint    `json:"tenth_year"`
	TenthMarks                   float64 `json:"tenth_marks"`
	TwelfthBoard                 string  `json:"twelfth_board"`
	TwelfthYear                  uint    `json:"twelfth_year"`
	TwelfthMarks                 float64 `json:"twelfth_marks"`
	EntranceExam                 string  `json:"entrance_exam"`
	EntranceExamRank             uint    `json:"entrance_exam_rank"`
	Category                     string  `json:"category"`
	CategoryRank                 uint    `json:"category_rank"`
	CurrentAddress               string  `json:"current_address"`
	PermanentAddress             string  `json:"permanent_address"`
	FriendName                   string  `json:"friend_name"`
	FriendPhone                  string  `json:"friend_phone"`
	Resume                       string  `json:"resume"`
	StatusName                   string  `json:"status_name"`
	Frozen                       bool    `json:"frozen"`
}

func getStudentsByRole(ctx *gin.Context) {
	pid, err := util.ParseUint(ctx.Param("pid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var applied []ApplicantsByRole
	fetchApplicantDetails(ctx, pid, &applied)

	var srids []uint
	for _, applicant := range applied {
		srids = append(srids, applicant.StudentRCID)
	}

	var allStudentsRC []rc.StudentRecruitmentCycle
	rc.FetchStudentBySRID(ctx, srids, &allStudentsRC)

	var allStudentsRCMap = make(map[uint]*rc.StudentRecruitmentCycle)
	for i := range allStudentsRC {
		allStudentsRCMap[allStudentsRC[i].ID] = &allStudentsRC[i]
	}

	var sid []uint
	for _, student := range allStudentsRC {
		sid = append(sid, student.StudentID)
	}

	var allStudents []student.Student
	student.FetchStudentsByID(ctx, sid, &allStudents)

	var allStudentsMap = make(map[uint]*student.Student)
	for i := range allStudents {
		allStudentsMap[allStudents[i].ID] = &allStudents[i]
	}

	var validApplicants []studentAdminsideResponse
	for _, s := range applied {
		// if allStudentsRCMap[student.StudentRCID].IsFrozen {
		// 	continue
		// }

		applicant_details := studentAdminsideResponse{}
		applicant_details.ID = s.StudentRCID
		applicant_details.Resume = s.ResumeLink
		applicant_details.StatusName = s.Name

		studentRC := allStudentsRCMap[s.StudentRCID]
		sid := allStudentsRCMap[s.StudentRCID].StudentID

		student := allStudentsMap[sid]

		applicant_details.Name = student.Name
		applicant_details.Email = student.IITKEmail

		applicant_details.CPI = studentRC.CPI
		applicant_details.ProgramDepartmentID = studentRC.ProgramDepartmentID
		applicant_details.SecondaryProgramDepartmentID = studentRC.SecondaryProgramDepartmentID

		applicant_details.TenthBoard = student.TenthBoard
		applicant_details.TenthYear = student.TenthYear
		applicant_details.TenthMarks = student.TenthMarks
		applicant_details.TwelfthBoard = student.TwelfthBoard
		applicant_details.TwelfthYear = student.TwelfthYear
		applicant_details.TwelfthMarks = student.TwelfthMarks
		applicant_details.EntranceExam = student.EntranceExam
		applicant_details.EntranceExamRank = student.EntranceExamRank
		applicant_details.Category = student.Category
		applicant_details.CategoryRank = student.CategoryRank
		applicant_details.CurrentAddress = student.CurrentAddress
		applicant_details.PermanentAddress = student.PermanentAddress
		applicant_details.FriendName = student.FriendName
		applicant_details.FriendPhone = student.FriendPhone
		applicant_details.Frozen = studentRC.IsFrozen

		validApplicants = append(validApplicants, applicant_details)
	}

	ctx.JSON(http.StatusOK, validApplicants)
}