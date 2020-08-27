package com.herbertgao.telegram.database.service;

import com.herbertgao.telegram.database.entity.ExamDate;
import com.herbertgao.telegram.database.entity.ExamDateExample;
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
     *
     * @return
     */
    public List<ExamDate> getExamList(LocalDateTime now, boolean desc) {
        ExamDateExample queryExample = new ExamDateExample();
        queryExample.createCriteria()
                .andExamYearBeginDateLessThanOrEqualTo(now)
                .andExamYearEndDateGreaterThanOrEqualTo(now)
                .andIsDeleteEqualTo(false);
        if (desc) {
            queryExample.setOrderByClause(" exam_year_begin_date desc ");
        }
        return mapper.selectByExample(queryExample);
    }

}
