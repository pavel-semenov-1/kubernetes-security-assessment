package com.ibm.demo.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

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
    @NoArgsConstructor
    public static class StudentId implements Serializable {
        private Long subjectId;
        private Long userId;
    }
}