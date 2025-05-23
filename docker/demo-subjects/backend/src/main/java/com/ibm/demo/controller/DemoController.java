package com.ibm.demo.controller;

import com.ibm.demo.entity.Student;
import com.ibm.demo.entity.Subject;
import com.ibm.demo.service.StudentService;
import com.ibm.demo.service.SubjectService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;

@RestController
@RequestMapping("/api/subjects")
public class DemoController {

    @Autowired
    private SubjectService subjectService;
    @Autowired
    private StudentService studentService;

    @GetMapping
    public List<Subject> getAllSubjects() {
        return subjectService.findAllSubjects();
    }

    @GetMapping("/{id}")
    public ResponseEntity<Subject> getSubjectById(@PathVariable Long id) {
        Optional<Subject> subject = subjectService.findSubjectById(id);
        return subject.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());
    }

    @PostMapping
    public Subject createSubject(@RequestBody Subject subject) {
        return subjectService.saveSubject(subject);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteSubject(@PathVariable Long id) {
        subjectService.deleteSubject(id);
        return ResponseEntity.noContent().build();
    }

    @PostMapping("/enroll")
    public Student enrollStudent(@RequestBody Student student) {
        return studentService.saveStudent(student);
    }

    @GetMapping("/{id}/students")
    public ResponseEntity<List<Student>> getStudentsBySubjectId(@PathVariable Long id) {
        List<Student> students = studentService.findStudentsBySubjectId(id);
        if (students.isEmpty()) {
            return ResponseEntity.noContent().build();
        }
        return ResponseEntity.ok(students);
    }

    @GetMapping("/user/{userId}")
    public ResponseEntity<List<Student>> getStudentsByUserId(@PathVariable Long userId) {
        List<Student> students = studentService.findStudentsByUserId(userId);
        if (students.isEmpty()) {
            return ResponseEntity.noContent().build();
        }
        return ResponseEntity.ok(students);
    }
}
