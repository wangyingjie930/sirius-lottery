package com.sirius.lottery.java.domain.repository;

import com.sirius.lottery.java.domain.entity.Instance;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface InstanceRepository extends JpaRepository<Instance, Long> {
    Optional<Instance> findByInstanceId(String instanceId);
}
