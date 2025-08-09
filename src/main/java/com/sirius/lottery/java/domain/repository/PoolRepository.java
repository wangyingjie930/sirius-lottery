package com.sirius.lottery.java.domain.repository;

import com.sirius.lottery.java.domain.entity.Pool;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface PoolRepository extends JpaRepository<Pool, Long> {
    List<Pool> findByInstanceId(String instanceId);
}
