package com.sirius.lottery.java.domain.repository;

import com.sirius.lottery.java.domain.entity.WinRecord;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface WinRecordRepository extends JpaRepository<WinRecord, Long> {
}
