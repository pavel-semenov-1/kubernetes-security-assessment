package com.ibm.demo.service;

import com.ibm.demo.entity.Student;
import com.ibm.demo.repository.StudentRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class StudentService {

    @Autowired
    private StudentRepository studentRepository;

    public List<Student> findAllStudents() {
        return studentRepository.findAll();
    }

    public Student saveStudent(Student student) {
        return studentRepository.save(student);
    }

    public List<Student> findStudentsBySubjectId(Long subjectId) {
        return studentRepository.findBySubjectId(subjectId);
    }

    public List<Student> findStudentsByUserId(Long studentId) {
        return studentRepository.findByUserId(studentId);
    }
}