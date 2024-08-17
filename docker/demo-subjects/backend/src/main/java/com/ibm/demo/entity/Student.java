package com.ibm.demo.entity;

import jakarta.persistence.*;
import lombok.Data;

import java.io.Serializable;

@Entity
@Data
@Table(name = "students")
@IdClass(Student.StudentId.class)
public class Student {
    @Id
    private Long subjectId;
    @Id
    private Long userId;

    @Data
    class StudentId implements Serializable {
        private Long subjectId;
        private Long userId;
    }
}