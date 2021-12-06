package com.herbertgao.telegram.database.service;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.herbertgao.telegram.database.entity.ExamDate;
import com.herbertgao.telegram.database.mapper.ExamDateMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

/**
 * @program: gaokao_bot
 * @description: ExamDateService
 * @author: HerbertGao
 * @create: 2019-06-09 01:37
 **/
@Service
public class ExamDateService {

    @Autowired
    private ExamDateMapper mapper;

    /**
     * @param now
     * @return
     */
    public List<ExamDate> getExamList(LocalDateTime now, boolean desc) {
        QueryWrapper<ExamDate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .le(ExamDate::getExamYearBeginDate, now)
                .ge(ExamDate::getExamYearEndDate, now)
                .eq(ExamDate::getIsDelete, false);
        if (desc) {
            wrapper.lambda().orderByDesc(ExamDate::getExamYearBeginDate);
        }
        return mapper.selectList(wrapper);
    }

    /**
     * @param year
     * @return
     */
    public List<ExamDate> getExamByYear(Integer year) {
        QueryWrapper<ExamDate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(ExamDate::getExamYear, year)
                .eq(ExamDate::getIsDelete, false);
        return mapper.selectList(wrapper);
    }
}
