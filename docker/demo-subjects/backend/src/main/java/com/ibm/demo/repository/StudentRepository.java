package com.ibm.demo.repository;

import com.ibm.demo.entity.Student;
import com.ibm.demo.entity.Subject;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface StudentRepository extends JpaRepository<Student, Student.StudentId> {
    List<Student> findBySubjectId(Long subjectId);
    List<Student> findByUserId(Long userId);
}