package com.sirius.lottery.java.domain.repository;

import com.sirius.lottery.java.domain.entity.Prize;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.Optional;

@Repository
public interface PrizeRepository extends JpaRepository<Prize, Long> {
    List<Prize> findByPoolId(Long poolId);

    @Query("SELECT p FROM Prize p JOIN Pool pool ON p.poolId = pool.id WHERE p.prizeId = :prizeId AND pool.instanceId = :instanceId")
    Optional<Prize> findByPrizeIdAndInstanceId(@Param("prizeId") String prizeId, @Param("instanceId") String instanceId);
}
